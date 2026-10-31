# Analyse approfondie : l'édition collaborative de documents

## Qu'est-ce qu'un éditeur collaboratif ?

Un éditeur collaboratif permet à plusieurs personnes -- et à des agents IA -- de modifier le même document en même temps, chaque modification apparaissant instantanément pour tous. Plus de « enregistrer et actualiser ». Plus de « ta version contre ma version ». Plus de conflits de fusion.

C'est la technologie qui fait tourner Google Docs, Notion et Figma. Polysee la propose sous la forme d'un composant intégrable que n'importe quel produit peut utiliser.

---

## L'architecture en un coup d'œil

```
┌─────────────────────────────────────────────────────┐
│                  Collaborative Editor                │
│                                                      │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐          │
│  │  User A   │  │  User B   │  │ AI Agent  │         │
│  │ (browser) │  │ (browser) │  │ (server)  │         │
│  └─────┬─────┘  └─────┬─────┘  └─────┬─────┘       │
│        │               │               │             │
│  ┌─────▼───────────────▼───────────────▼─────┐      │
│  │        Hocuspocus Sync Server              │      │
│  │   (WebSocket, conflict resolution)         │      │
│  └─────────────────────┬─────────────────────┘      │
│                        │                             │
│  ┌─────────────────────▼─────────────────────┐      │
│  │         Yjs (CRDT Document Model)          │      │
│  │   Automatic conflict-free merging          │      │
│  └─────────────────────┬─────────────────────┘      │
│                        │                             │
│  ┌─────────────────────▼─────────────────────┐      │
│  │           PostgreSQL Storage               │      │
│  │   Binary document state + snapshots        │      │
│  └───────────────────────────────────────────┘      │
└─────────────────────────────────────────────────────┘
```

---

## Comment fonctionne la synchronisation en temps réel

### Le problème

Quand deux personnes modifient le même paragraphe au même moment, laquelle des deux modifications « l'emporte » ? Les systèmes traditionnels recourent soit au verrouillage (« quelqu'un d'autre modifie cette section »), soit à la règle du dernier enregistrement (qui écrase des modifications sans prévenir). Les deux approches sont mauvaises.

### La solution : les CRDT

Polysee s'appuie sur **Yjs**, une bibliothèque CRDT (Conflict-free Replicated Data Type, type de données répliqué sans conflit). Voici son fonctionnement, en termes simples :

1. L'éditeur de chaque utilisateur conserve une copie locale du document
2. Chaque frappe génère une **opération** (insérer un caractère à la position X, supprimer un caractère à la position Y)
3. Les opérations sont envoyées au **serveur de synchronisation Hocuspocus** via WebSocket
4. Le serveur diffuse les opérations à tous les autres éditeurs connectés
5. L'algorithme CRDT de Yjs **garantit** que tous les éditeurs convergent vers le même état du document, quel que soit l'ordre d'arrivée des opérations

**Aucun conflit. Aucune perte de données. Aucun verrouillage.** Deux utilisateurs peuvent taper dans le même mot au même moment : le résultat reste toujours cohérent.

### Gestion de la connexion

- **Connexion WebSocket** au serveur Hocuspocus (port 3001)
- **Reconnexion automatique** en cas de déconnexion
- **Suivi de l'état** : connexion en cours, connecté, déconnecté
- **Gestion des délais d'attente** avec des valeurs par défaut configurables (5 à 10 secondes)

---

## Ce que vous pouvez rédiger

L'éditeur repose sur **Tiptap** (une surcouche moderne de ProseMirror) et prend en charge un contenu riche :

| Type de contenu | Description |
|---|---|
| **Texte enrichi** | Gras, italique, souligné, barré |
| **Titres** | Hiérarchie H1, H2, H3 |
| **Listes** | Listes à puces et listes numérotées/ordonnées |
| **Blocs de code** | Coloration syntaxique pour plus de 50 langages de programmation, avec sélecteur de langage |
| **Code en ligne** | Mise en forme à chasse fixe avec surlignage personnalisé |
| **Citations** | Blocs de citation en retrait |
| **Tableaux** | Prise en charge complète des tableaux, avec colonnes redimensionnables |
| **Images** | Images intégrées par téléversement, glisser-déposer ou collage depuis le presse-papiers |
| **Pièces jointes** | Fichiers PDF, Word, Excel et Markdown, avec affichage des métadonnées |
| **Séparateurs horizontaux** | Séparateurs visuels entre sections |
| **Liens** | URL cliquables avec un style personnalisé |

### Commandes slash

Tapez `/` n'importe où dans le document pour ouvrir un menu de commandes :

| Commande | Action |
|---|---|
| `/image` | Téléverser et intégrer une image |
| `/file` | Joindre un document (PDF, Word, Excel, Markdown) |
| `/h1`, `/h2`, `/h3` | Insérer un titre |
| `/bullet` | Commencer une liste à puces |
| `/numbered` | Commencer une liste numérotée |
| `/code` | Insérer un bloc de code |
| `/quote` | Insérer une citation |
| `/bold`, `/italic` | Appliquer une mise en forme au texte |

Les commandes se filtrent en temps réel pendant la saisie : `/co` affiche immédiatement « Bloc de code ».

### Collage intelligent

- **Collez du Markdown** : il est converti automatiquement en texte enrichi
- **Collez des images** depuis le presse-papiers -- elles sont téléversées et intégrées dans le texte
- **Déposez des fichiers** sur l'éditeur -- les images sont intégrées, les documents joints

---

## Présence et conscience de l'activité

Lorsque plusieurs personnes se trouvent dans le même document, chacune voit qui est là et ce que font les autres :

### Utilisateurs actifs

Un indicateur de présence, dans le coin supérieur, affiche :
- Les avatars de tous les utilisateurs connectés
- Le nombre de personnes en train d'éditer
- Un dédoublonnage des utilisateurs ouverts dans plusieurs onglets

### Curseurs en direct

Le curseur de chaque utilisateur apparaît dans le document en temps réel :
- **Codé par couleur**, pour distinguer les utilisateurs
- **Étiqueté** avec le nom de l'utilisateur
- **Curseurs IA** dotés d'une animation pulsée distinctive, pour différencier les modifications d'un bot de celles d'un humain

Cette visibilité évite le classique « tu es encore en train de modifier ? » et rend la collaboration naturelle.

---

## Historique des versions

Chaque modification significative d'un document est conservée sous forme d'**instantané** (snapshot) :

### Instantanés automatiques

Le système crée automatiquement des instantanés sur la base de :
- **Hachage du contenu** (SHA-256) : seul un contenu réellement différent déclenche un nouvel instantané
- **Suivi du nombre de mots** : chaque instantané enregistre la longueur du document
- **Attribution à l'auteur** : qui a effectué la modification (utilisateur ou agent IA, avec nom et avatar)

### Panneau d'historique

Une frise chronologique visuelle permet aux utilisateurs de :
1. **Parcourir** toutes les versions passées, avec horodatage et informations sur l'auteur
2. **Prévisualiser** n'importe quel instantané sans toucher au document actuel
3. **Restaurer** n'importe quelle version en un clic

La restauration crée un nouvel instantané au point de restauration : vous ne perdez jamais l'état actuel -- vous pouvez toujours « annuler l'annulation ».

### Stockage des instantanés

Chaque instantané contient :
- L'état binaire complet du document (format Yjs)
- Une empreinte du contenu pour le dédoublonnage
- Le nombre de mots
- Les informations sur l'auteur (utilisateur ou agent, avec avatar)
- L'horodatage

---

## Édition assistée par IA

C'est là que l'éditeur collaboratif prend toute sa puissance. Les agents IA peuvent modifier des documents par programme, avec la même synchronisation en temps réel que les utilisateurs humains.

### Insertion en flux continu

L'IA peut écrire du contenu dans le document en temps réel :

```
AI Agent                        Document
   │                               │
   │── start ─────────────────────>│  (cursor appears)
   │── position(paragraph 3) ─────>│  (cursor moves)
   │── chunk("The system") ───────>│  (text appears)
   │── chunk(" provides") ────────>│  (text continues)
   │── chunk(" real-time") ───────>│  (text continues)
   │── complete ──────────────────>│  (cursor disappears)
   │                               │
```

Les utilisateurs regardent l'IA écrire en temps réel, comme s'ils regardaient un collègue taper. Ils peuvent continuer à modifier d'autres parties du document en parallèle.

### Application de diffs

L'IA peut proposer des modifications structurées du contenu existant :

| Format de diff | Description |
|---|---|
| **Markdown Replace** | Remplace tout le contenu du document par du nouveau Markdown |
| **Unified Diff** | Insère, supprime ou remplace des sections précises |
| **JSON Patch** | Correctifs structurés pour des modifications précises |

Chaque application de diff enregistre :
- Qui l'a appliqué (utilisateur ou agent)
- Le type de source (ai-bot, user, system)
- Une description de ce qui a changé
- L'horodatage

### Client bot

Pour l'édition par IA côté serveur, un client bot sans interface se connecte au document :
1. Il génère un jeton JWT pour s'authentifier
2. Il ouvre une connexion Yjs vers le serveur Hocuspocus
3. Il applique les modifications par programme
4. Il apparaît comme un agent nommé dans la liste de présence

Les agents IA peuvent ainsi modifier des documents même lorsqu'aucun humain ne les a ouverts -- pratique pour des tâches de fond comme la génération de documentation ou la migration de contenu.

---

## Gestion des fichiers et des images

### Téléversement d'images

Trois façons d'ajouter des images :
1. **Commande slash** (`/image`) -- ouvre le sélecteur de fichiers
2. **Glisser-déposer** -- déposez des fichiers image sur l'éditeur
3. **Collage depuis le presse-papiers** -- collez directement des captures d'écran

Les images sont redimensionnées automatiquement à une largeur maximale de 800 px et passent par le même flux de téléversement S3 sécurisé en 3 étapes que les pièces jointes du chat.

### Pièces jointes

Prise en charge des fichiers autres que les images :
- PDF, Word (.docx), Excel (.xlsx), Markdown (.md)
- Taille maximale : 10 Mo par fichier
- Les fichiers s'affichent sous forme de cartes indiquant le nom, le type et la taille
- Un clic suffit pour les télécharger

### Zone de dépôt

Lorsque vous faites glisser des fichiers au-dessus de l'éditeur, une surcouche visuelle apparaît pour confirmer la zone de dépôt. Cela évite les dépôts accidentels et offre un retour clair.

---

## Export et conversion

Les documents peuvent être extraits dans plusieurs formats :

| Format | Cas d'usage |
|---|---|
| **Texte brut** | Indexation pour la recherche, affichage simple |
| **Markdown** | Documentation pour développeurs, génération de sites statiques |
| **HTML** | Rendu web, contenu d'e-mails |
| **JSON** | Accès par programme, schéma ProseMirror |

La conversion Markdown bidirectionnelle vous permet de :
- Importer du contenu Markdown existant dans l'éditeur
- Exporter les documents collaboratifs en Markdown pour les versionner
- Faire des allers-retours de contenu entre l'éditeur et d'autres systèmes

---

## Contrôle d'accès

L'accès aux documents est contrôlé à plusieurs niveaux :

| Niveau | Contrôle |
|---|---|
| **Authentification** | Sessions, jetons d'accès personnels, clés d'API, JWT (pour les bots) |
| **Périmètre du projet** | Les documents sont rattachés à des projets ; être membre d'un projet donne accès à ses documents |
| **Basé sur les rôles** | Les propriétaires et les mainteneurs ont un accès en lecture/écriture ; les autres rôles en lecture seule |
| **Protection IDOR** | La propriété de la fonctionnalité est vérifiée par rapport au projet avant d'accorder l'accès |

---

## Gestion des erreurs et fiabilité

### Error Boundary de l'éditeur

Un Error Boundary React enveloppe le composant de l'éditeur. Si l'éditeur plante (contenu malformé, erreur d'extension, etc.), une interface de repli soignée s'affiche au lieu d'une page blanche. Les utilisateurs peuvent recharger sans perdre de données (le serveur détient le dernier état).

### Reprise de connexion

- Reconnexion automatique en cas de coupure WebSocket
- Synchronisation de l'état du document à la reconnexion (Yjs se charge de la fusion)
- Nettoyage des données de curseur et de présence à la déconnexion (pas de curseurs fantômes)

### Sécurité des données

- Toutes les modifications sont enregistrées dans PostgreSQL au format binaire Yjs
- Les instantanés permettent une restauration à un instant donné
- Le hachage du contenu empêche les instantanés en double de gonfler le stockage

---

## REST API

L'éditeur est entièrement accessible par API pour les flux de travail automatisés :

| Méthode | Point de terminaison | Rôle |
|---|---|---|
| `GET` | `/api/features/{id}/documentation` | Récupérer le contenu actuel en Markdown |
| `PATCH` | `/api/features/{id}/documentation` | Remplacer le contenu par du nouveau Markdown |
| `POST` | `/api/features/{id}/documentation/apply-diff` | Appliquer un diff structuré |
| `POST` | `/api/features/{id}/documentation/stream-insert` | Lancer une insertion IA en flux continu (SSE) |
| `POST` | `/api/features/{id}/documentation/snapshot` | Créer un instantané manuel |
| `GET` | `/api/features/{id}/documentation/history` | Lister tous les instantanés |
| `GET` | `/api/features/{id}/documentation/history/{snapshotId}` | Obtenir le détail d'un instantané avec aperçu |
| `POST` | `/api/features/{id}/documentation/restore` | Revenir à un instantané précédent |

---

## Questions fréquentes

### Qu'est-ce que l'édition collaborative ?

L'édition collaborative, c'est lorsque plusieurs personnes travaillent sur le même document en même temps et voient instantanément les modifications des autres. Chaque frappe est synchronisée avec toutes les copies ouvertes, et les modifications faites au même moment sont fusionnées automatiquement : il n'y a jamais de « ta version contre ma version ».

### Plusieurs utilisateurs et des agents IA peuvent-ils modifier le même document en même temps ?

Oui. Un agent IA rejoint le document comme n'importe quel participant : il a son propre curseur nommé, apparaît dans la liste de présence, et ses modifications se synchronisent comme celles d'une personne. Les utilisateurs peuvent continuer à taper dans d'autres parties du document pendant que l'agent écrit, et chaque modification est attribuée à son auteur, humain ou agent.

### Comment sont gérés les conflits quand deux personnes modifient la même ligne ?

Le document est un CRDT (Yjs). Les modifications simultanées sont fusionnées de façon déterministe sur chaque client, si bien que toutes les copies finissent identiques, sans verrous, sans bouton « enregistrer » ni fusion manuelle. Les modifications faites hors ligne sont fusionnées de la même façon au retour de la connexion.

### Comment ajouter des curseurs en direct et la présence à un éditeur collaboratif ?

Ils sont fournis avec l'éditeur : chaque participant diffuse son curseur, sa sélection, son nom et sa couleur sur la même connexion de synchronisation, et l'éditeur les affiche pour tous les autres. Les curseurs des personnes qui partent sont nettoyés automatiquement, et les curseurs IA sont dessinés différemment pour que les modifications des bots se remarquent.

### Peut-on l'utiliser pour une documentation Markdown collaborative ?

Oui. Le Markdown peut être collé, importé, exporté et écrit via la REST API, et l'éditeur conserve un historique complet des versions avec restauration en un clic. Il convient donc à la documentation produit et aux spécifications qu'une équipe et ses agents IA maintiennent ensemble.

### Puis-je intégrer l'éditeur dans mon propre produit ?

Oui. Il est conçu comme un composant intégrable, avec son propre serveur de synchronisation, son stockage, son contrôle d'accès et sa REST API : un produit peut ainsi proposer une édition à la Google Docs sans construire lui-même la couche temps réel.

---

## En résumé

L'éditeur collaboratif n'est pas une simple zone de texte. C'est une **plateforme documentaire en temps réel** qui :

1. **Permet l'édition simultanée**, avec résolution automatique des conflits grâce aux CRDT
2. **Montre qui est où**, avec des curseurs en direct et des indicateurs de présence
3. **Préserve l'historique**, avec des instantanés automatiques et une restauration en un clic
4. **Fait participer l'IA** comme collaborateur en temps réel -- qui écrit, modifie et propose des changements
5. **Prend en charge un contenu riche** -- texte, code, tableaux, images et pièces jointes
6. **Exporte partout** -- Markdown, HTML, JSON, texte brut
7. **Reste fiable** grâce aux Error Boundaries, à la reprise de connexion et au stockage persistant

C'est l'infrastructure documentaire sur laquelle reposent des produits comme Notion, Confluence et Google Docs -- disponible sous forme de composant intégrable pour n'importe quelle application.
