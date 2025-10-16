# Merchant Dialog Tree

```mermaid
flowchart TD
    %% Starting Node
    START[start: Well met, traveler! I am Gareth...]
    
    %% Primary Branches from Start
    START --> SHOW_WARES[show_wares: What goods for sale?]
    START --> ROAD_NEWS[road_news: Any news from road?]
    START --> BANDIT_TROUBLE[bandit_trouble: Heard about bandit trouble]
    START --> POLITE_FAREWELL[polite_farewell: Just passing through]
    
    %% Wares Branch
    SHOW_WARES --> OPEN_SHOP[open_shop: Browse inventory]
    SHOW_WARES --> BEST_WEAPON[best_weapon: What's your best weapon?]
    SHOW_WARES --> HAGGLE_PRICES[haggle_prices: Can you lower prices?]
    SHOW_WARES --> OTHER_TOPICS[other_topics: Ask about other things]
    
    %% Shop Branch
    OPEN_SHOP --> POLITE_FAREWELL
    
    %% Road News Branch
    ROAD_NEWS --> STRANGE_CREATURES[strange_creatures: Tell about creatures]
    ROAD_NEWS --> MISSING_CARAVANS[missing_caravans: Missing caravans?]
    ROAD_NEWS --> OFFER_BANDIT_HELP[offer_bandit_help: I can help with bandits]
    
    %% Bandit Trouble Branch
    BANDIT_TROUBLE --> OFFER_BANDIT_HELP
    BANDIT_TROUBLE --> STRANGE_CREATURES
    BANDIT_TROUBLE --> OTHER_TOPICS
    
    %% Help Offer Branch
    OFFER_BANDIT_HELP --> ACCEPT_BANDIT_QUEST[accept_bandit_quest: Accept the mission]
    OFFER_BANDIT_HELP --> NEGOTIATE_PAYMENT[negotiate_payment: What's the reward?]
    OFFER_BANDIT_HELP --> OTHER_TOPICS
    
    %% Quest Acceptance Branch
    ACCEPT_BANDIT_QUEST --> QUEST_ACCEPTED_END[quest_accepted_end: Quest taken]
    
    %% Negotiation Branch
    NEGOTIATE_PAYMENT --> ACCEPT_BANDIT_QUEST
    NEGOTIATE_PAYMENT --> DEMAND_MORE_PAYMENT[demand_more_payment: Need more gold]
    NEGOTIATE_PAYMENT --> OTHER_TOPICS
    
    %% Demand More Branch
    DEMAND_MORE_PAYMENT --> ACCEPT_BANDIT_QUEST
    DEMAND_MORE_PAYMENT --> POLITE_FAREWELL
    
    %% Creatures Branch
    STRANGE_CREATURES --> OFFER_BANDIT_HELP
    STRANGE_CREATURES --> OTHER_TOPICS
    
    %% Missing Caravans Branch
    MISSING_CARAVANS --> OFFER_BANDIT_HELP
    MISSING_CARAVANS --> OTHER_TOPICS
    
    %% Other Topics Branch
    OTHER_TOPICS --> SHOW_WARES
    OTHER_TOPICS --> ROAD_NEWS
    OTHER_TOPICS --> POLITE_FAREWELL
    
    %% Ending Nodes (Auto "End dialog" choice)
    POLITE_FAREWELL:::endNode
    QUEST_ACCEPTED_END:::endNode
    OPEN_SHOP:::endNode
    
    %% Styling
    classDef startNode fill:#4CAF50,stroke:#2E7D32,stroke-width:3px,color:#fff
    classDef endNode fill:#F44336,stroke:#C62828,stroke-width:3px,color:#fff
    classDef tradeNode fill:#FF9800,stroke:#E65100,stroke-width:2px,color:#fff
    classDef newsNode fill:#2196F3,stroke:#1565C0,stroke-width:2px,color:#fff
    classDef questNode fill:#9C27B0,stroke:#6A1B9A,stroke-width:2px,color:#fff
    classDef negotiateNode fill:#795548,stroke:#3E2723,stroke-width:2px,color:#fff
    
    %% Apply styles
    class START startNode
    class SHOW_WARES,OPEN_SHOP,BEST_WEAPON,HAGGLE_PRICES tradeNode
    class ROAD_NEWS,STRANGE_CREATURES,MISSING_CARAVANS newsNode
    class OFFER_BANDIT_HELP,ACCEPT_BANDIT_QUEST questNode
    class NEGOTIATE_PAYMENT,DEMAND_MORE_PAYMENT negotiateNode
```

## Legend
- 🟢 **Start Node**: Entry point (start)
- 🔴 **End Nodes**: Auto "End dialog" choice appears
- 🟠 **Trade Nodes**: Shopping and commerce
- 🔵 **News Nodes**: Road information and rumors
- 🟣 **Quest Nodes**: Mission offers and acceptance
- 🟤 **Negotiation Nodes**: Payment discussions

## Key Dialog Paths

### 🛒 **Shopping Path**
```
start → show_wares → open_shop → [END]
```

### 🗡️ **Quest Path (Direct)**
```
start → bandit_trouble → offer_bandit_help → accept_bandit_quest → quest_accepted_end
```

### 💰 **Negotiation Path**
```
start → road_news → offer_bandit_help → negotiate_payment → accept_bandit_quest → quest_accepted_end
```

### 📰 **Information Gathering Path**
```
start → road_news → strange_creatures → other_topics → show_wares
```

### 🤝 **Polite Exit Path**
```
start → polite_farewell → [END]
```

## Variables Affected
- `merchant_met`: Set to true at start
- `merchant_reputation`: Increases with positive interactions (0-3)
- `bought_items`: Set when shopping (placeholder)
- `heard_road_news`: Set when discussing road conditions
- `bandit_quest_from_merchant`: Set when accepting quest from merchant

## Special Features
- **Multiple Quest Entry Points**: Can reach bandit quest through different conversation paths
- **Circular Navigation**: `other_topics` allows returning to main conversation branches
- **Commerce Integration**: `open_shop` node ready for inventory system integration
- **Reputation System**: Merchant remembers player interactions through reputation variable