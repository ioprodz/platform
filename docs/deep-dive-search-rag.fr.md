# Analyse approfondie : recherche et RAG (génération augmentée par la recherche)

## Qu'est-ce que le RAG ?

Les grands modèles de langage sont puissants, mais ils ont une faiblesse majeure : ils ne connaissent que ce qui figurait dans leurs données d'entraînement. Ils ignorent tout des documents de votre entreprise, des fonctionnalités de votre produit et des échanges de votre équipe.

La **génération augmentée par la recherche (RAG, Retrieval-Augmented Generation)** répond à ce problème en donnant à l'IA accès à vos données au moment de la requête :

```
Without RAG:
  User: "What are our authentication requirements?"
  AI: "Authentication typically involves..." (generic, possibly wrong)

With RAG:
  User: "What are our authentication requirements?"
  System: [searches your docs, finds the auth spec]
  AI: "According to your Feature-Auth-2024 spec, authentication requires
       OAuth2 with PKCE flow, session timeout of 30 minutes, and MFA
       for admin roles." (specific, grounded, accurate)
```

Le moteur de recherche et RAG rend cela possible : il indexe votre contenu, en comprend le sens et retrouve les passages les plus pertinents dès que l'IA en a besoin.

---

## L'architecture en un coup d'œil

```
┌──────────────────────────────────────────────────────┐
│                   Search & RAG Engine                 │
│                                                       │
│  ┌────────────────────────────────────────────────┐  │
│  │              Content Sources                    │  │
│  │  Documents | Features | Chats | Tickets | PRs  │  │
│  └───────────────────┬────────────────────────────┘  │
│                      │ (events)                       │
│  ┌───────────────────▼────────────────────────────┐  │
│  │           Indexing Pipeline                      │  │
│  │  Extract → Hash → Embed → Store                 │  │
│  └───────────────────┬────────────────────────────┘  │
│                      │                                │
│  ┌───────────────────▼────────────────────────────┐  │
│  │        PostgreSQL + pgvector                     │  │
│  │  Full-text index | Vector embeddings             │  │
│  └───────────────────┬────────────────────────────┘  │
│                      │                                │
│  ┌───────────────────▼────────────────────────────┐  │
│  │           Hybrid Search                         │  │
│  │  (0.7 × vector score) + (0.3 × text score)     │  │
│  └───────────────────┬────────────────────────────┘  │
│                      │                                │
│  ┌──────────┐  ┌─────▼──────┐  ┌──────────────────┐ │
│  │  Search   │  │    RAG     │  │   AI Agent       │ │
│  │   API     │  │  Retrieve  │  │  Tool Calling    │ │
│  └──────────┘  └────────────┘  └──────────────────┘ │
└──────────────────────────────────────────────────────┘
```

---

## Fonctionnement, étape par étape

### Étape 1 : extraction du contenu

Dès qu'un contenu est créé ou modifié n'importe où dans la plateforme, des **extracteurs** le convertissent en texte indexable. Chaque type d'entité dispose de son propre extracteur :

| Type d'entité | Ce qui est extrait |
|---|---|
| **Documents** | Texte intégral du document issu de l'éditeur collaboratif (binaire Yjs → texte brut) |
| **Fonctionnalités** | Nom, statut, description et contenu complet du document |
| **Cas d'utilisation** | Nom, acteurs, critères d'acceptation |
| **Pull requests** | Titre, état, type, noms de branches, métadonnées |
| **Évaluations** | Nom, description, questions et revues d'itération |
| **Messages de chat** | Contenu du fil (jusqu'à 50 messages par fil) |
| **Tickets de support** | Titre, statut, description, métadonnées |

La recherche couvre donc tout : pas seulement les documents, mais aussi les conversations, les revues de code, les tickets et bien plus encore.

### Étape 2 : détection des modifications

Avant de générer un nouvel embedding (une opération payante), le système vérifie si le contenu a réellement changé :

1. Calcul d'une **empreinte SHA-256** du texte extrait
2. Comparaison avec l'empreinte enregistrée pour cette entité
3. **Si elles sont identiques** : pas de nouvel embedding (économie de coûts)
4. **Si elles diffèrent** : génération d'un nouvel embedding et mise à jour de l'empreinte

Vous ne payez ainsi la génération d'embeddings que lorsque le contenu change vraiment.

### Étape 3 : génération des embeddings

Le texte extrait est envoyé à un **modèle d'embedding** qui le convertit en vecteur de grande dimension, c'est-à-dire une représentation numérique de son sens :

```
"OAuth2 authentication with PKCE flow"
    → [0.023, -0.156, 0.891, 0.045, ..., -0.234]  (1536 dimensions)
```

**Modèles d'embedding pris en charge :**

| Modèle | Dimensions | Idéal pour |
|---|---|---|
| `text-embedding-3-small` | 1536 | Usage général au meilleur coût |
| `text-embedding-3-large` | 3072 | Précision maximale |
| `text-embedding-ada-002` | 1536 | Compatibilité avec l'existant |

Le modèle se configure par projet : vous arbitrez entre coût et qualité selon vos besoins.

### Étape 4 : stockage

Le texte d'origine et le vecteur d'embedding sont tous deux stockés dans **PostgreSQL avec pgvector** :

- **Index plein texte** (`tsvector`) : pour la correspondance par mots-clés
- **Colonne vectorielle** : pour la recherche par similarité sémantique
- **Métadonnées** : type d'entité, identifiant d'entité, empreinte du contenu, modèle d'embedding, horodatage

Utiliser PostgreSQL à la fois pour les données relationnelles et pour la recherche vectorielle évite d'administrer une base de données supplémentaire : tout est au même endroit.

### Étape 5 : recherche

À chaque requête, le système lance une **recherche hybride** qui combine deux stratégies distinctes :

#### Recherche vectorielle (sémantique)
- La requête est convertie en embedding avec le même modèle
- pgvector trouve les vecteurs les plus proches selon la distance cosinus
- On retrouve ainsi les contenus qui **veulent dire** la même chose, même formulés avec d'autres mots
- Exemple : une recherche sur « sécurité de connexion » remonte des contenus portant sur les « mécanismes d'authentification »

#### Recherche plein texte (mots-clés)
- Les fonctions natives de PostgreSQL `tsvector` et `plainto_tsquery` trouvent les correspondances de mots exactes et partielles
- On capte ainsi les contenus qui emploient les **termes exacts** saisis par l'utilisateur
- Exemple : une recherche sur « OAuth2 » remonte les contenus contenant précisément ce terme

#### Score combiné

```
final_score = (vector_score × 0.7) + (text_score × 0.3)
```

La pondération 70/30 privilégie la compréhension sémantique tout en récompensant les correspondances exactes de mots-clés. Ces pondérations sont configurables pour chaque requête.

**Pourquoi l'hybride ?** La recherche purement sémantique rate les termes exacts (une recherche sur « JIRA-1234 » doit trouver précisément ce ticket). La recherche purement par mots-clés rate le sens (une recherche sur « authentification » doit trouver les contenus qui parlent de « connexion »). L'approche hybride réunit le meilleur des deux.

---

## Indexation pilotée par les événements

Le contenu est indexé **automatiquement**, dès qu'il change. Pas de tâches cron. Pas de déclenchement manuel. Pas de résultats périmés.

Le système écoute les événements de chaque module de contenu :

| Module | Événements |
|---|---|
| **Documents** | Création, mise à jour, suppression |
| **Fonctionnalités** | Création, mise à jour du nom/de la description/du statut, archivage, suppression |
| **Cas d'utilisation** | Création, mise à jour, suppression |
| **Pull requests** | Création, mise à jour, suppression |
| **Évaluations** | Création, mise à jour, suppression |
| **Messages de chat** | Création, mise à jour, suppression |
| **Tickets de support** | Création, mise à jour, suppression |

Lorsqu'un événement se déclenche :
1. L'extracteur approprié génère le texte indexable
2. Le pipeline d'indexation vérifie si un nouvel embedding est nécessaire
3. Le cas échéant, un nouvel embedding est généré et stocké
4. L'index de recherche est immédiatement à jour

Une **file d'attente asynchrone** à concurrence configurable (3 par défaut) traite les tâches d'indexation sans bloquer l'application principale.

---

## Récupération RAG pour les agents IA

La fonction de récupération RAG fait le lien entre la recherche et l'IA :

```
ragRetrieve(projectId, query, topK = 5)
  → [
      { content: "...", source: "Feature: Auth System", score: 0.92 },
      { content: "...", source: "Use Case: User Login", score: 0.87 },
      { content: "...", source: "Chat: Security Discussion", score: 0.71 },
    ]
```

Lorsqu'un agent IA a besoin d'une information :
1. Il appelle l'outil `search_features` avec une requête en langage naturel
2. L'outil exécute une recherche hybride sur l'index du projet
3. Les meilleurs résultats sont renvoyés avec leur contenu, leur source et leur score de pertinence
4. L'agent intègre ces résultats à sa réponse en **citant ses sources**

C'est ce qui empêche les hallucinations. L'agent ne devine pas : il récupère le contenu réel de votre base de connaissances et s'y réfère.

---

## Modes de recherche

Le système propose deux modes de fonctionnement :

### Mode texte seul

- Utilise uniquement la recherche plein texte de PostgreSQL
- Aucune génération d'embeddings (coût d'IA nul pour l'indexation)
- Correspondance par mots-clés uniquement
- Idéal pour : les petits projets, les environnements soucieux des coûts, ou lorsque la recherche sémantique n'est pas nécessaire

### Mode hybride

- Combine recherche vectorielle et recherche plein texte
- Nécessite un fournisseur d'embeddings (OpenAI)
- Correspondance sémantique et par mots-clés
- Idéal pour : tout projet dont les utilisateurs doivent trouver des contenus par leur sens, et pas seulement par les mots exacts

Les projets peuvent changer de mode à tout moment. Lors du passage en mode hybride, une **tâche en arrière-plan** génère automatiquement les embeddings de tous les enregistrements existants en mode texte seul.

---

## Configuration

### Paramètres par projet

Chaque projet peut configurer :
- **Mode de recherche** : text_only ou hybrid
- **Fournisseur d'embeddings** : le fournisseur OpenAI à utiliser
- **Modèle d'embedding** : le modèle servant à générer les embeddings
- **Héritage** : utiliser les valeurs globales par défaut ou les remplacer pour le projet

### Valeurs globales par défaut

Les administrateurs peuvent définir des valeurs par défaut à l'échelle du système, dont héritent les projets ; le comportement reste ainsi cohérent, sauf remplacement explicite.

### Réindexation

Une réindexation manuelle peut être lancée :
- Lors d'un changement de modèle d'embedding (les anciens embeddings doivent être régénérés)
- Lors du passage du mode texte seul au mode hybride
- En cas de doute sur l'intégrité des données

Le processus de réindexation :
1. Extrait le contenu de toutes les entités du projet
2. Traite les entités par lots de 10
3. Suit la progression (nombre total d'entités, entités indexées)
4. Signale les erreurs sans interrompre le lot
5. Met à jour le statut de la configuration une fois terminé

---

## Résultats de recherche

Chaque résultat de recherche comprend :

| Champ | Description |
|---|---|
| **title** | Nom de l'entité ou titre déduit |
| **snippet** | Extrait pertinent avec le contexte autour de la correspondance |
| **score** | Score de pertinence combiné (de 0 à 1) |
| **vectorScore** | Composante de similarité sémantique |
| **textScore** | Composante de correspondance par mots-clés |
| **entityType** | Nature du contenu (fonctionnalité, cas d'utilisation, chat, etc.) |
| **url** | Lien direct vers l'entité source |
| **metadata** | Détails propres à l'entité (statut, responsable, etc.) |

Les résultats renvoient directement au contenu source : l'utilisateur passe d'un clic du résultat au document, à la conversation ou au ticket concerné.

---

## Référence de l'API

### Points de terminaison de recherche

| Méthode | Point de terminaison | Rôle |
|---|---|---|
| `POST` | `/api/projects/{id}/search` | Exécuter une requête de recherche hybride |
| `PUT` | `/api/projects/{id}/search` | Lancer une réindexation complète |
| `GET` | `/api/projects/{id}/search` | Obtenir l'état actuel de l'index |

### Points de terminaison de configuration

| Méthode | Point de terminaison | Rôle |
|---|---|---|
| `GET` | `/api/projects/{id}/search/config` | Obtenir les paramètres de recherche du projet |
| `POST` | `/api/projects/{id}/search/config` | Modifier les paramètres de recherche du projet |
| `GET` | `/api/admin/search/embedding-config` | Obtenir les valeurs globales par défaut |
| `POST` | `/api/admin/search/embedding-config` | Définir les valeurs globales par défaut |
| `GET` | `/api/admin/search/reindex-status` | Suivre la progression de la réindexation |
| `GET` | `/api/admin/search/debug` | Diagnostiquer l'état de l'index de recherche |

---

## Intégration côté client

Un hook React assure une intégration fluide de la recherche :

```
useSearch({ projectId, debounceMs: 300, limit: 20 })
```

Fonctionnalités :
- **Requêtes temporisées (debounce)** : la recherche part 300 ms après la dernière frappe de l'utilisateur
- **Annulation des requêtes** : la recherche précédente est annulée dès qu'une nouvelle requête arrive
- **Mise en cache des résultats** : 30 secondes avant péremption, 60 secondes avant le nettoyage du cache
- **Préchargement** : chargement anticipé des résultats pour les requêtes prévisibles

---

## Choix techniques

### Pourquoi PostgreSQL + pgvector (plutôt qu'une base vectorielle dédiée) ?

- **Une seule base de données** : aucune infrastructure supplémentaire à administrer (ni Pinecone, ni Weaviate, ni Qdrant)
- **Cohérence transactionnelle** : les embeddings sont mis à jour dans la même transaction que le contenu
- **Écosystème mature** : la recherche plein texte de PostgreSQL a fait ses preuves à toutes les échelles
- **Coût** : aucun abonnement SaaS supplémentaire pour une base vectorielle

### Pourquoi la recherche hybride (plutôt que purement vectorielle) ?

- **Les correspondances exactes comptent** : une recherche sur « PROJ-1234 » exige une correspondance par mots-clés, pas une similarité sémantique
- **Acronymes et codes** : les termes propres à un domaine nécessitent souvent une correspondance exacte
- **Équilibre configurable** : la pondération 70/30 peut être ajustée pour chaque requête selon les cas d'usage

### Pourquoi une indexation pilotée par les événements (plutôt que par lots) ?

- **Fraîcheur en temps réel** : le contenu est consultable par la recherche dès sa création
- **Maîtrise des coûts** : seul le contenu modifié reçoit un nouvel embedding
- **Aucune tâche planifiée** : pas de complexité liée à cron, pas d'attente « jusqu'à la prochaine indexation »

---

## En résumé

Le moteur de recherche et RAG n'est pas un simple champ de recherche. C'est une **couche de connaissances intelligente** qui :

1. **Indexe tout** : documents, conversations, tickets, revues de code, évaluations
2. **Comprend le sens** grâce aux embeddings vectoriels sémantiques
3. **Reste à jour** grâce à une indexation en temps réel pilotée par les événements
4. **Ancre les réponses de l'IA** dans vos données réelles via l'API de récupération RAG
5. **Équilibre précision et rappel** grâce à une recherche hybride configurable
6. **Passe à l'échelle simplement** avec PostgreSQL + pgvector (sans infrastructure de base de données supplémentaire)
7. **Respecte les budgets** grâce aux empreintes de contenu, qui évitent de régénérer inutilement des embeddings

Voilà comment rendre l'IA fiable en production : en lui donnant accès à vos données réelles, et pas seulement à son corpus d'entraînement.
