# Town Elder Dialog Tree

```mermaid
flowchart TD
    %% Starting Node
    START[start: Welcome to our humble village...]
    
    %% Primary Branches from Start
    START --> ASK_PROPHECY[ask_prophecy: Ancient prophecy knowledge]
    START --> ASK_VILLAGE[ask_village: Tell me about village]
    START --> ASK_WORK[ask_work: Looking for work]
    START --> END_FRIENDLY[end_friendly: Farewell for now]
    
    %% Prophecy Branch
    ASK_PROPHECY --> CLAIM_HERO[claim_hero: I am the foretold hero]
    ASK_PROPHECY --> CASUAL_INTEREST[casual_interest: Just curious]
    ASK_PROPHECY --> ASK_DARKNESS[ask_darkness: What darkness?]
    
    %% Village Information Branch
    ASK_VILLAGE --> ASK_DANGERS[ask_dangers: What dangers exist?]
    ASK_VILLAGE --> ASK_TRADE[ask_trade: About trade/economy]
    ASK_VILLAGE --> END_FRIENDLY
    
    %% Work Branch
    ASK_WORK --> ACCEPT_BANDIT_QUEST[accept_bandit_quest: Accept bandit mission]
    ASK_WORK --> ASK_PATROLS[ask_patrols: Ask about patrols]
    ASK_WORK --> GENERAL_INFO[general_info: General village info]
    
    %% Hero Claim Branch
    CLAIM_HERO --> ASK_TRIALS[ask_trials: What trials await?]
    CLAIM_HERO --> ACCEPT_CHALLENGE[accept_challenge: Ready for trials]
    
    %% Casual Interest Branch
    CASUAL_INTEREST --> ASK_EVENTS[ask_events: Recent strange events]
    CASUAL_INTEREST --> ASK_BELIEF[ask_belief: Do you believe prophecy?]
    
    %% Darkness Investigation Branch
    ASK_DARKNESS --> VOLUNTEER_INVESTIGATE[volunteer_investigate: I'll investigate]
    ASK_DARKNESS --> ASK_PROTECTION[ask_protection: Village defense]
    
    %% Dangers Branch
    ASK_DANGERS --> OFFER_HELP[offer_help: I can help]
    
    %% Trade Branch
    ASK_TRADE --> ACCEPT_BANDIT_QUEST
    ASK_TRADE --> GENERAL_INFO
    ASK_TRADE --> END_FRIENDLY
    
    %% Patrols Branch
    ASK_PATROLS --> VOLUNTEER_INVESTIGATE
    ASK_PATROLS --> ASK_DARKNESS
    ASK_PATROLS --> ASK_PROTECTION
    
    %% General Info Branch
    GENERAL_INFO --> ASK_DANGERS
    GENERAL_INFO --> ASK_EVENTS
    GENERAL_INFO --> END_FRIENDLY
    
    %% Trials Branch
    ASK_TRIALS --> ACCEPT_CHALLENGE
    ASK_TRIALS --> QUEST_GUIDANCE[quest_guidance: Where to find trials?]
    ASK_TRIALS --> END_PREPARATION[end_preparation: Need more preparation]
    
    %% Challenge Accept Branch
    ACCEPT_CHALLENGE --> END_QUEST_GIVEN[end_quest_given: Quest accepted]
    ACCEPT_CHALLENGE --> QUEST_GUIDANCE
    
    %% Events Branch
    ASK_EVENTS --> ASK_BELIEF
    ASK_EVENTS --> VOLUNTEER_INVESTIGATE
    ASK_EVENTS --> ASK_PROTECTION
    
    %% Belief Branch
    ASK_BELIEF --> CLAIM_HERO
    ASK_BELIEF --> OFFER_HELP
    ASK_BELIEF --> END_HELPFUL[end_helpful: Time will tell]
    
    %% Protection Branch
    ASK_PROTECTION --> OFFER_HELP
    ASK_PROTECTION --> VOLUNTEER_INVESTIGATE
    ASK_PROTECTION --> GENERAL_INFO
    
    %% Help Offers Branch
    OFFER_HELP --> QUEST_GUIDANCE
    VOLUNTEER_INVESTIGATE --> QUEST_GUIDANCE
    
    %% Bandit Quest Branch
    ACCEPT_BANDIT_QUEST --> END_QUEST_GIVEN
    ACCEPT_BANDIT_QUEST --> ASK_DANGERS
    
    %% Ending Nodes (Auto "End dialog" choice)
    END_FRIENDLY:::endNode
    END_HELPFUL:::endNode
    END_QUEST_GIVEN:::endNode
    END_PREPARATION:::endNode
    
    %% Styling
    classDef startNode fill:#4CAF50,stroke:#2E7D32,stroke-width:3px,color:#fff
    classDef endNode fill:#F44336,stroke:#C62828,stroke-width:3px,color:#fff
    classDef questNode fill:#FF9800,stroke:#E65100,stroke-width:2px,color:#fff
    classDef infoNode fill:#2196F3,stroke:#1565C0,stroke-width:2px,color:#fff
    classDef actionNode fill:#9C27B0,stroke:#6A1B9A,stroke-width:2px,color:#fff
    
    %% Apply styles
    class START startNode
    class ACCEPT_BANDIT_QUEST,ACCEPT_CHALLENGE,VOLUNTEER_INVESTIGATE questNode
    class ASK_PROPHECY,ASK_VILLAGE,ASK_WORK,GENERAL_INFO,ASK_TRIALS,ASK_EVENTS infoNode
    class CLAIM_HERO,OFFER_HELP,ASK_BELIEF actionNode
```

## Legend
- 🟢 **Start Node**: Entry point (start)
- 🔴 **End Nodes**: Auto "End dialog" choice appears
- 🟠 **Quest Nodes**: Accept missions/challenges
- 🔵 **Info Nodes**: Gather information
- 🟣 **Action Nodes**: Character development choices

## Key Dialog Paths

### 🗡️ **Hero's Journey Path**
```
start → ask_prophecy → claim_hero → ask_trials → accept_challenge → end_quest_given
```

### 🏹 **Bandit Quest Path**
```
start → ask_work → accept_bandit_quest → end_quest_given
```

### 🛡️ **Investigation Path**
```
start → ask_village → ask_dangers → offer_help → quest_guidance
```

### 📖 **Lore Discovery Path**
```
start → ask_prophecy → casual_interest → ask_events → ask_belief
```

## Variables Affected
- `met_elder`: Set to true at start
- `elder_trust`: Increases based on player choices (0-4)
- `knows_prophecy`: Set when discussing prophecy
- `bandit_quest_accepted`: Set when accepting bandit mission
- `trial_quest_given`: Set when accepting trial challenge
- `knows_trials`: Set when learning about trials