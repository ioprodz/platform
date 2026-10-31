# Analyse approfondie : le moteur d'IA

## Qu'est-ce qu'un moteur d'IA ?

Un moteur d'IA est la couche qui s'intercale entre votre produit et les grands modèles de langage (LLM). Sans lui, intégrer l'IA oblige à écrire du code spécifique pour chaque fournisseur, à gérer les clés d'API, le streaming, l'orchestration des outils et le suivi des coûts – le tout en partant de zéro.

Le moteur d'IA réunit tout cela dans une plateforme unique et configurable. Vous définissez des **agents** (qui est l'IA), vous leur donnez des **outils** (ce que l'IA peut faire) et vous les reliez à des **fournisseurs** (le modèle qui les fait tourner). Tout le reste – streaming, nouvelles tentatives, suivi des coûts, multi-tenant – est pris en charge pour vous.

---

## L'architecture en un coup d'œil

```
┌─────────────────────────────────────────────────────┐
│                   Your Application                  │
├─────────────────────────────────────────────────────┤
│                                                     │
│   ┌──────────┐   ┌──────────┐   ┌──────────┐        │
│   │  Agent A │   │  Agent B │   │  Agent C │        │
│   │ "Helper" │   │ "Analyst"│   │ "Coder"  │        │
│   └────┬─────┘   └─────┬────┘   └──────┬───┘        │
│        │               │               │            │
│   ┌────▼───────────────▼───────────────▼────┐       │
│   │          Tool Calling Loop              │       │
│   │  (reason → call tool → observe → repeat)│       │
│   └────┬────────────────────────────────┬───┘       │
│        │                                │           │
│   ┌────▼──────┐                  ┌──────▼─────┐     │
│   │ Built-in  │                  │     MCP    │     │
│   │   Tools   │                  │ (external) │     │
│   └───────────┘                  └────────────┘     │
│                                                     │
├─────────────────────────────────────────────────────┤
│            Provider Abstraction Layer               │
│   ┌──────────┐  ┌──────────┐  ┌──────────┐          │
│   │  OpenAI  │  │  Claude  │  │  Ollama  │          │
│   │  (cloud) │  │  (cloud) │  │  (local) │          │
│   └──────────┘  └──────────┘  └──────────┘          │
├─────────────────────────────────────────────────────┤
│   Usage Tracking  │  Cost Management  │  Analytics  │
└─────────────────────────────────────────────────────┘
```

---

## Concepts clés

### 1. Les fournisseurs

Un **fournisseur** est une connexion à un service de LLM. Chaque fournisseur comprend :
- Un point d'accès d'API et des identifiants
- La liste des modèles disponibles
- Des règles de tarification pour estimer les coûts

**Fournisseurs pris en charge :**

| Fournisseur | Type | Modèles | Appel d'outils | Streaming |
|---|---|---|---|---|
| **OpenAI** | Cloud | GPT-4o, GPT-4, GPT-3.5, etc. | Oui | Oui |
| **Anthropic** | Cloud | Claude 3.5 Sonnet, Claude 3, etc. | Oui | Oui |
| **Ollama** | Local | Mistral 7B, Llama, etc. | Limité | Oui |

Grâce à cette couche d'abstraction, le code de votre application ne fait jamais référence à un modèle précis. Vous configurez le fournisseur utilisé par chaque agent, et le moteur se charge de traduire votre requête dans le format d'API propre à chaque fournisseur.

**Pourquoi c'est important :** si un nouveau modèle sort demain, vous changez une seule valeur de configuration. Aucune modification de code. Aucun redéploiement.

### 2. Les agents

Un **agent** est un personnage d'IA nommé, doté de :
- **Un nom et un avatar** : la façon dont il apparaît dans l'interface
- **Un prompt système** : les instructions qui définissent son comportement et sa personnalité
- **Un fournisseur attribué** : le LLM qui le fait tourner (il peut remplacer celui défini par défaut pour le projet)
- **Un accès aux outils** : les outils qu'il est autorisé à utiliser
- **Des connexions à des serveurs MCP** : les serveurs d'outils externes qu'il peut joindre

Les agents sont rattachés à un projet : chaque projet peut donc disposer de ses propres agents, configurés différemment.

**Exemple :** un agent « Analyste produit » peut s'appuyer sur Claude pour un raisonnement nuancé, avoir accès aux outils de recherche de fonctionnalités et de documentation, et avoir pour consigne de toujours citer ses sources. Un agent « Assistant rapide » peut utiliser un modèle Ollama local pour répondre vite, et sans frais, aux questions simples.

### 3. L'appel d'outils

C'est l'appel d'outils qui rend les agents d'IA utiles au-delà de la simple conversation. Au lieu de se contenter de générer du texte, un agent peut :

1. **Raisonner** sur les informations dont il a besoin
2. **Appeler un outil** pour récupérer des données ou effectuer une action
3. **Observer** le résultat renvoyé par l'outil
4. **Recommencer** jusqu'à disposer de suffisamment d'informations pour répondre

Cette boucle peut s'exécuter jusqu'à 50 fois par requête, ce qui permet aux agents de mener seuls des workflows complexes en plusieurs étapes.

**Parmi les outils intégrés :**
- `search_features` : rechercher du contenu pertinent dans la base de connaissances
- `send_message` : envoyer des messages à des utilisateurs ou à des canaux
- `label_activity` : étiqueter et classer le travail de l'agent lui-même
- `get_tool_guide` : obtenir des instructions détaillées pour les usages d'outils complexes

### 4. MCP (Model Context Protocol)

MCP est un protocole standard qui relie les agents d'IA à des serveurs d'outils externes. Voyez-le comme « l'USB de l'IA » : tout système qui parle MCP peut instantanément doter un agent de nouvelles capacités.

**Fonctionnement :**
1. Un serveur MCP expose un ensemble d'outils via JSON-RPC 2.0
2. L'agent découvre les outils disponibles au moment de la connexion
3. Pendant la conversation, l'agent appelle ces outils exactement comme les outils intégrés
4. Les résultats reviennent dans la boucle de raisonnement de l'agent

**Transports pris en charge :**
- **stdio** : communication entre processus locaux (pour les outils exécutés sur le même serveur)
- **HTTP** : communication avec un serveur distant (pour les outils hébergés ailleurs)

**Pourquoi c'est important :** vous étendez les capacités des agents sans toucher au moteur d'IA lui-même. Un agent doit interroger votre CRM ? Déployez un serveur MCP qui encapsule l'API de votre CRM. Il doit vérifier les stocks ? Déployez un autre serveur MCP. L'agent découvre et utilise automatiquement les nouveaux outils.

---

## Déclencheurs d'activation

Les agents ne restent pas inactifs. Ils se déclenchent en réponse à des événements :

| Déclencheur | Fonctionnement |
|---|---|
| **Mention dans le chat** | Un utilisateur mentionne un agent avec @ dans le chat. L'agent reçoit le contexte de la conversation et répond. |
| **Réponse à un questionnaire** | Un utilisateur remplit une série de questions. L'agent traite les réponses et agit en conséquence. |
| **Événement du calendrier** | Un événement planifié se déclenche. L'agent exécute son workflow à l'heure configurée. |
| **Fin d'une tâche de code** | Un agent de code externe termine une tâche. L'agent d'IA traite les résultats et en rend compte. |

Grâce à cette conception événementielle, les agents s'intègrent naturellement aux workflows existants. Les utilisateurs n'ont pas à apprendre une nouvelle interface : il leur suffit de mentionner un agent dans la conversation qu'ils ont déjà en cours.

---

## Capacités vocales

Le moteur d'IA intègre le traitement de la voix :

- **Reconnaissance vocale (Speech-to-Text)** : les messages vocaux sont automatiquement transcrits en arrière-plan. La transcription est stockée avec l'audio, ce qui rend le contenu vocal consultable par recherche et lisible par l'IA.
- **Synthèse vocale (Text-to-Speech)** : lorsqu'un utilisateur envoie un message vocal, l'agent peut répondre par un audio généré. Le prompt système s'adapte automatiquement au contexte vocal (réponses concises, pas de mise en forme Markdown, durée de 30 à 60 secondes).

---

## Suivi de l'utilisation et maîtrise des coûts

Chaque requête d'IA fait l'objet d'un suivi :

| Indicateur | Description |
|---|---|
| **Tokens en entrée** | Nombre de tokens envoyés au modèle |
| **Tokens en sortie** | Nombre de tokens générés par le modèle |
| **Coût estimé** | Calculé à partir de règles de tarification configurables par fournisseur et par modèle |
| **Latence** | Temps de réponse de bout en bout |
| **Statut** | Succès, erreur, délai dépassé ou limite de débit atteinte |
| **Type de fonctionnalité** | La capacité à l'origine de la requête (chat, génération de titre, prompt système, etc.) |
| **Attribution à l'utilisateur** | L'utilisateur qui a lancé la requête |

Ces données permettent :
- **Un contrôle budgétaire** par projet
- **La comparaison des fournisseurs** (coût, qualité, rapidité)
- **Des rapports d'utilisation** pour la facturation et la planification des capacités
- **La détection d'anomalies** en cas d'emballement des coûts

---

## Suivi des activités

Chaque interaction d'un agent est consignée dans un **enregistrement d'activité** qui contient :
- Le prompt système complet et le contexte de la conversation
- Tous les appels d'outils, avec leurs entrées, leurs sorties et leur numéro d'itération
- La durée et le statut d'achèvement
- Le contexte de la fonctionnalité ou du ticket associé

Vous disposez ainsi d'une piste d'audit complète de ce que l'IA a fait, pourquoi et en combien de temps – un élément essentiel pour le débogage, la conformité et la confiance.

---

## Multi-tenant

Le moteur d'IA est conçu dès le départ pour les produits multi-tenant :

- **Agents rattachés au projet** : chaque projet possède ses propres agents, configurés indépendamment
- **Partage des fournisseurs** : les administrateurs peuvent créer des fournisseurs au niveau du système et en accorder l'accès à des projets précis
- **Isolation de l'utilisation** : le nombre de tokens et les coûts sont suivis par projet
- **Héritage de la configuration** : les projets peuvent utiliser les valeurs par défaut globales ou les remplacer

---

## Points d'intégration

Le moteur d'IA se connecte aux autres composants de la plateforme :

| Composant | Intégration |
|---|---|
| **Chat** | Les agents reçoivent les messages, diffusent leurs réponses en streaming et utilisent des éléments interactifs |
| **Recherche et RAG** | Les agents appellent `search_features` pour ancrer leurs réponses dans des données réelles |
| **Éditeur collaboratif** | Les agents modifient les documents par programmation, via l'insertion en streaming et l'application de diffs |
| **Calendrier** | Des déclencheurs planifiés lancent les workflows des agents |
| **Stockage de fichiers** | Les enregistrements vocaux et les pièces jointes transitent par S3 |

---

## En résumé

Le moteur d'IA n'est pas un chatbot. C'est un environnement d'exécution complet pour agents qui :

1. **Fait abstraction des fournisseurs**, pour ne jamais vous enfermer dans un seul modèle
2. **Orchestre les outils**, pour que les agents puissent agir concrètement
3. **S'étend grâce à MCP**, pour enrichir les capacités sans modifier le code
4. **Suit tout**, pour que vous compreniez les coûts, les performances et le comportement
5. **Isole les tenants**, pour que chaque projet fonctionne de manière indépendante

C'est la couche d'infrastructure d'IA dont tout produit moderne a besoin – et la construire de zéro demanderait des mois de travail d'ingénierie dédiée.
