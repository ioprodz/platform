# Analyse approfondie : chat et collaboration en temps réel

## Qu'est-ce qu'un système de chat nativement IA ?

Les systèmes de chat traditionnels font circuler du texte entre des humains. Un système de chat nativement IA traite les **agents IA comme des participants à part entière** : ils lisent les messages, comprennent le contexte, diffusent leurs réponses en continu, appellent des outils et affichent des éléments interactifs, le tout dans le même fil de conversation que celui des humains.

Il ne s'agit pas d'un « chat IA » distinct greffé sur votre produit. C'est une couche de communication unifiée où humains et IA collaborent naturellement.

---

## L'architecture en un coup d'œil

```
┌──────────────────────────────────────────────────────┐
│                     Chat System                       │
│                                                       │
│  ┌─────────────┐  ┌─────────────┐  ┌──────────────┐ │
│  │   Humans    │  │  AI Agents  │  │ CLI Agents   │ │
│  │  (users)    │  │ (assistants)│  │  (devices)   │ │
│  └──────┬──────┘  └──────┬──────┘  └──────┬───────┘ │
│         │                │                 │          │
│  ┌──────▼─────────────────▼─────────────────▼──────┐ │
│  │              Message Stream                      │ │
│  │  text | voice | files | interactive elements     │ │
│  └──────────────────┬──────────────────────────────┘ │
│                     │                                 │
│  ┌──────────────────▼──────────────────────────────┐ │
│  │           Real-Time Delivery (SSE)               │ │
│  │   streaming responses | presence | notifications │ │
│  └──────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────┘
```

---

## Types de conversations

Le système de chat prend en charge quatre modèles de conversation, qui couvrent tous les schémas de communication :

| Type | Description | Exemple |
|---|---|---|
| **Privée (utilisateur)** | Messagerie individuelle entre deux utilisateurs | Message direct à un collègue |
| **Privée (agent)** | Conversation individuelle avec un agent IA | Poser une question à l'agent « Product Analyst » |
| **Groupe** | Conversations à plusieurs participants | Canal d'une équipe projet réunissant utilisateurs et agents |
| **Contextuelle** | Chat intégré à une fonctionnalité ou une page précise | Fil de discussion sur la spécification d'une fonctionnalité |

Tous les types de conversation reposent sur la même infrastructure : messages, pièces jointes, diffusion en temps réel et participation de l'IA fonctionnent de façon identique, quel que soit le type de conversation.

---

## Types et contenu des messages

### Messages texte

Messages texte standard, avec :
- Rendu Markdown complet (titres, listes, blocs de code, liens, gras, italique)
- @-mentions avec références cliquables vers les utilisateurs et les agents
- Copie dans le presse-papiers pour n'importe quel message

### Messages vocaux

Un parcours complet de messagerie vocale :

```
Record → Preview → Send → Auto-Transcribe
```

1. **Enregistrement** : capture audio dans le navigateur, avec visualisation du volume en temps réel et minuteur
2. **Aperçu** : réécoute avant l'envoi, avec possibilité de supprimer et de réenregistrer
3. **Envoi** : l'audio est téléversé vers un stockage cloud et le message est délivré avec des commandes de lecture
4. **Transcription** : une conversion parole-texte en arrière-plan transforme automatiquement l'audio en texte consultable par recherche

Les messages vocaux sont particulièrement utiles pour :
- Les utilisateurs « mobile first » qui ont besoin de communiquer les mains libres
- Les équipes de terrain qui font leurs comptes rendus depuis les chantiers
- L'accessibilité : les utilisateurs pour qui la saisie au clavier est difficile
- Saisir les nuances que le texte fait perdre

### Pièces jointes

Glisser-déposer ou clic pour téléverser, avec :
- **Images** : JPEG, PNG, GIF, WebP, avec affichage intégré et superposition au survol
- **Documents** : PDF, TXT, Markdown, avec des fiches indiquant le nom, le type et la taille du fichier
- **Plusieurs fichiers** par message
- **Bandeau d'aperçu** avant l'envoi : vérifiez et retirez des fichiers avant qu'ils ne partent
- **Bouton de téléchargement** sur les fichiers reçus

Le téléversement suit un parcours sécurisé en trois étapes :
1. Demander au serveur une URL de téléversement présignée
2. Téléverser directement vers le stockage cloud (S3)
3. Confirmer le téléversement et associer les métadonnées au message

Grâce à cette conception, les fichiers ne transitent jamais par le serveur applicatif, ce qui réduit la charge et la latence.

### Messages système

Notifications automatiques pour les événements de la conversation :
- Création d'un groupe (« Alice a créé ce groupe »)
- Changements de participants (« Bob a été ajouté à la conversation »)
- Événements système personnalisés issus des intégrations

---

## Éléments interactifs

C'est ici que le système de chat dépasse la messagerie traditionnelle. Les messages peuvent contenir des **composants actionnables** avec lesquels les utilisateurs interagissent directement :

| Type d'élément | Description | Exemple |
|---|---|---|
| **Proposition de diff** | L'IA suggère une modification de code ou de document. L'utilisateur peut l'accepter ou la rejeter sur place. | L'agent propose un correctif : « Remplacer X par Y à la ligne 42 » ; l'utilisateur clique sur Accepter |
| **Insertion en continu** | L'IA écrit du contenu dans un document en temps réel, avec des commandes pour valider ou rejeter | L'agent rédige une section de spécification pendant que l'utilisateur regarde |
| **Formulaire** | Collecte de données structurées au fil du chat | L'agent demande des paramètres de configuration via un formulaire |
| **Confirmation** | Point de décision oui/non, accompagné de son contexte | L'agent demande : « Dois-je déployer ceci en staging ? » |
| **Série de questions** | Questionnaire à plusieurs questions, avec cases à cocher et boutons radio | Questionnaire d'accueil avec relances pilotées par l'IA |

Chaque élément interactif suit son propre état (`pending` → `accepted`/`rejected`/`completed`/`failed`), et chaque changement d'état déclenche des événements auxquels d'autres systèmes peuvent réagir.

**Pourquoi c'est important :** les utilisateurs n'ont pas besoin de quitter la conversation pour agir. L'IA propose, l'utilisateur décide et le système exécute, le tout au même endroit.

---

## Mentions et étiquetage

Le système de @-mentions relie les humains et l'IA :

1. **Tapez `@`** dans le champ de saisie pour ouvrir le menu de suggestions
2. **Filtrez** en tapant : les suggestions se resserrent en temps réel
3. **Sélectionnez** un utilisateur, un agent IA ou un appareil CLI
4. **La mention est intégrée** au message sous forme d'étiquette structurée

Lorsqu'un message mentionne un agent IA, le système, automatiquement :
- Achemine le message vers la chaîne de traitement de l'agent
- Joint l'historique récent de la conversation comme contexte
- Commence à diffuser en continu la réponse de l'agent dans le chat

### Agents persistants

Pour les échanges suivis avec un agent donné, la fonction « agent persistant » (sticky agent) inclut automatiquement l'agent dans chaque message, sans qu'il soit nécessaire de le @-mentionner. Ce réglage est enregistré par conversation et se conserve d'une session à l'autre.

---

## Diffusion en temps réel

Les messages sont délivrés en temps réel grâce aux **Server-Sent Events (SSE)** :

```
Client                          Server
  │                                │
  │── GET /api/chat/events ───────>│
  │                                │
  │<── event: message ────────────│  (new message)
  │<── event: message ────────────│  (streaming chunk)
  │<── event: heartbeat ──────────│  (keep-alive)
  │<── event: message ────────────│  (another message)
  │                                │
  │── (connection timeout) ───────>│
  │── GET /api/chat/events ───────>│  (auto-reconnect)
  │                                │
```

**Une interrogation intelligente qui s'adapte au contexte :**
- **Pendant la diffusion d'une réponse IA** : interrogation toutes les secondes pour une arrivée fluide des tokens
- **En fonctionnement normal** : interrogation toutes les 3 secondes, par souci d'efficacité
- **Sur mobile** : durée de connexion maximale plus courte (2 minutes) pour préserver la batterie
- **Expiration de la connexion** : 10 minutes au maximum, avec reconnexion automatique

**La détection des changements** est optimisée :
- Détection des mises à jour de contenu pour les réponses diffusées en continu
- Détection des changements de contexte pour les mises à jour de métadonnées
- Détection des nouveaux messages pour la diffusion standard
- Seules les données modifiées sont transmises

---

## Réponses IA en continu

Lorsqu'un agent IA répond, la réponse s'affiche dans le chat token par token :

1. Un message provisoire apparaît immédiatement, avec l'avatar de l'agent
2. Le texte s'affiche progressivement, au rythme de la génération par le modèle
3. Les appels d'outils sont affichés en temps réel (les utilisateurs voient ce que fait l'agent)
4. Le message est finalisé à la fin de la génération

Un **tampon de flux** (stream buffer) gère ce flux :
- Les fragments sont collectés et vidés toutes les 200 ms
- Des écritures atomiques en base de données évitent les conditions de concurrence
- Des écritures séquentielles garanties préservent la cohérence des messages

Il en résulte une expérience naturelle et réactive : les utilisateurs voient l'IA « réfléchir » au lieu d'attendre une réponse complète.

---

## Notifications et suivi des messages non lus

### Compteurs de non-lus

Chaque conversation enregistre, pour chaque utilisateur, un **horodatage de dernière lecture**. Les compteurs de non-lus sont calculés en comparant cet horodatage aux nouveaux messages, hors messages système.

### Notifications sonores

Alertes sonores configurables pour les nouveaux messages, avec gestion des préférences de chaque utilisateur.

### Événements de notification en temps réel

Un canal SSE dédié délivre les événements de notification (aperçus des nouveaux messages, informations sur l'expéditeur) afin que l'interface mette à jour les badges et indicateurs sans interroger les listes de conversations.

---

## Gestion des discussions de groupe

Gestion complète du cycle de vie des conversations à plusieurs participants :

| Opération | Description |
|---|---|
| **Créer** | Choisir les participants (utilisateurs + agents), nommer le groupe ; le créateur est ajouté automatiquement |
| **Inviter** | Ajouter de nouveaux participants, avec notification par message système |
| **Quitter** | Se retirer, avec boîte de dialogue de confirmation |
| **Modifier** | Mettre à jour le nom du groupe et la liste des participants |

Les groupes offrent les mêmes fonctionnalités que les discussions privées : pièces jointes, messages vocaux, agents IA, éléments interactifs et mentions.

---

## Pagination et historique

L'historique du chat utilise une **pagination par curseur** pour un chargement efficace :

- **Chargement initial** : les 10 messages les plus récents
- **Charger plus** : 15 messages par page, déclenché par le défilement près du haut
- **Bidirectionnelle** : chargement des messages plus anciens ou plus récents par rapport à un curseur
- **Gestion du défilement** : défilement automatique vers le bas à l'arrivée de nouveaux messages, badge « Nouveaux messages » lorsque l'utilisateur est remonté dans l'historique

---

## Statistiques d'utilisation

Le système de chat intègre des statistiques :

| Indicateur | Description |
|---|---|
| **Total des messages** | Volume global de messages |
| **Conversations uniques** | Nombre de fils actifs |
| **Utilisateurs actifs** | Participants humains distincts |
| **Messages par membre** | Nombre de messages par utilisateur |
| **Messages par agent** | Nombre de messages par agent |
| **Répartition quotidienne** | Activité des membres et des agents dans le temps |
| **Répartition par projet** | Volume de messages par projet |

Ces indicateurs permettent de :
- Comprendre les modes de communication de l'équipe
- Mesurer l'adoption et l'efficacité des agents IA
- Repérer les zones de forte activité susceptibles de nécessiter davantage de ressources
- Produire des rapports pour la direction et les parties prenantes

---

## Sécurité et multi-tenant

| Fonctionnalité | Description |
|---|---|
| **Isolation par projet** | Les messages et conversations sont cloisonnés par projet |
| **Authentification** | Par session (NextAuth), jetons d'accès personnels (Personal Access Tokens), clés d'API |
| **Autorisation** | Les opérations d'écriture exigent une authentification ; l'accès en lecture est configurable |
| **Piste d'audit** | Chaque message porte l'attribution de son expéditeur, des horodatages et des métadonnées de contexte |

---

## Contexte des messages et extensibilité

Chaque message transporte un **objet de contexte** qui peut inclure :

- Les références aux utilisateurs étiquetés
- Des indicateurs d'état de diffusion en continu
- Les métadonnées des agents CLI (nom de l'appareil, modèle, temps de réponse)
- Des références de réponse (base des fils de discussion)
- Des tableaux d'éléments interactifs
- Les métadonnées des pièces jointes
- Les métadonnées vocales (type MIME, durée, état de la transcription)
- Les informations de fuseau horaire
- L'état d'indexation pour la recherche

Grâce à ce modèle de contexte extensible, de nouvelles fonctionnalités peuvent être ajoutées aux messages sans migration de schéma : les champs de contexte sont stockés en JSON et interprétés par la couche d'interface.

---

## En résumé

Le système de chat ne se limite pas à la messagerie. C'est un **environnement d'exécution de la collaboration** qui :

1. **Réunit humains et IA** dans les mêmes fils de conversation
2. **Prend en charge des médias riches** : texte, voix, fichiers et composants interactifs
3. **Diffuse en temps réel**, avec un streaming adaptatif et des notifications
4. **Permet d'agir** grâce à des éléments interactifs qui laissent les utilisateurs décider et exécuter sans quitter le chat
5. **Assure un suivi complet** grâce aux statistiques, aux pistes d'audit et à la gestion des non-lus
6. **Passe à l'échelle sur plusieurs projets**, avec une isolation multi-tenant et des accès configurables

C'est la couche de communication dont ont besoin les produits modernes augmentés par l'IA : chaque conversation y devient une occasion d'automatisation intelligente.
