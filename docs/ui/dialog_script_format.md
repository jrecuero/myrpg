# Dialog Script Format

This document defines the JSON format for dialog scripts used by the Dialog Widget system. The format supports branching conversations, conditional logic, character definitions, and game event integration.

## File Structure

Dialog scripts are stored in the `assets/dialogs/` directory as JSON files:

```
assets/dialogs/
├── characters.json          # Character definitions
├── town_elder.json         # Sample NPC dialog
├── merchant_intro.json     # Merchant conversations
└── quest_giver.json       # Quest-related dialogs
```

## Character Definitions (`characters.json`)

Defines characters that appear in dialogs with their portrait information:

> **📚 Detailed Guide**: For comprehensive instructions on adding, removing, and managing characters, see [Dialog Character Management Guide](dialog_character_management.md)

```json
{
  "characters": {
    "town_elder": {
      "name": "Elder Aldric",
      "portrait": "elder_portrait.png",
      "voice_style": "wise",
      "description": "The wise elder of Rivertown"
    },
    "merchant": {
      "name": "Trader Gareth",
      "portrait": "merchant_portrait.png", 
      "voice_style": "cheerful",
      "description": "A traveling merchant"
    },
    "player": {
      "name": "Hero",
      "portrait": "player_portrait.png",
      "voice_style": "determined",
      "description": "The game protagonist"
    }
  }
}
```

## Dialog Script Structure

Each dialog file contains one or more conversation trees:

```json
{
  "dialog_id": "town_elder_intro",
  "title": "Meeting the Elder",
  "description": "First conversation with the town elder",
  "variables": {
    "met_elder": false,
    "elder_trust": 0,
    "knows_prophecy": false
  },
  "nodes": {
    "start": {
      "speaker": "town_elder",
      "text": "Welcome to our humble village, traveler. I am Elder Aldric.",
      "conditions": [],
      "actions": [
        {"type": "set_variable", "name": "met_elder", "value": true}
      ],
      "choices": [
        {
          "text": "Greetings, Elder. I seek information about the prophecy.",
          "target": "ask_prophecy",
          "conditions": [
            {"type": "variable", "name": "knows_prophecy", "operator": "equals", "value": false}
          ]
        },
        {
          "text": "Hello. Can you tell me about this village?",
          "target": "ask_village"
        },
        {
          "text": "Farewell for now.",
          "target": "end_friendly"
        }
      ]
    },
    "ask_prophecy": {
      "speaker": "town_elder", 
      "text": "Ah, the ancient prophecy... Few know of it these days. What brings you to seek such knowledge?",
      "conditions": [],
      "actions": [
        {"type": "set_variable", "name": "elder_trust", "value": 1}
      ],
      "choices": [
        {
          "text": "I believe I may be the one foretold in the prophecy.",
          "target": "claim_hero",
          "conditions": [
            {"type": "level", "operator": "greater_than", "value": 5}
          ]
        },
        {
          "text": "I'm just curious about local legends.",
          "target": "casual_interest"
        }
      ]
    },
    "ask_village": {
      "speaker": "town_elder",
      "text": "Rivertown has stood for three centuries. We are simple folk who work the land and trade with passing merchants.",
      "choices": [
        {
          "text": "Are there any dangers I should know about?",
          "target": "ask_dangers"
        },
        {
          "text": "Thank you for the information.",
          "target": "end_friendly"
        }
      ]
    },
    "end_friendly": {
      "speaker": "town_elder",
      "text": "May the winds guide your path, traveler. You are always welcome in Rivertown.",
      "actions": [
        {"type": "set_variable", "name": "elder_trust", "value": 2}
      ],
      "end": true
    }
  }
}
```

## Node Structure

Each dialog node contains the following fields:

### Required Fields
- **`speaker`**: Character ID from characters.json
- **`text`**: The dialog text to display

### Optional Fields
- **`conditions`**: Array of conditions that must be met to show this node
- **`actions`**: Array of actions to execute when this node is reached
- **`choices`**: Array of player response choices
- **`end`**: Boolean, if true, automatically adds an "End dialog" choice to close the conversation
- **`next`**: Auto-advance to another node (if no choices)

## Dialog Ending Behavior

### Automatic End Choice
When a dialog node has `"end": true` and no `choices` array, the system automatically adds an "End dialog" choice that closes the dialog when selected:

```json
{
  "farewell": {
    "speaker": "town_elder",
    "text": "May fortune favor your journey, traveler.",
    "end": true
  }
}
```

**Result**: After the text displays, shows a single "End dialog" button that closes the dialog when clicked.

### Manual End Choices
You can also create manual end choices by using empty `target` fields or pointing to nodes with `"end": true`:

```json
{
  "main_conversation": {
    "speaker": "merchant",
    "text": "What brings you to my shop?",
    "choices": [
      {
        "text": "Show me your wares.",
        "target": "show_items"
      },
      {
        "text": "Farewell.",
        "target": "goodbye"
      }
    ]
  },
  "goodbye": {
    "speaker": "merchant", 
    "text": "Safe travels!",
    "end": true
  }
}
```

### End Choice Priority
- If a node has both `choices` and `"end": true`, the existing choices take precedence
- Automatic "End dialog" choice only appears when `choices` array is empty or missing
- This allows for complex conversation trees with multiple ending paths

## Condition System

Conditions determine when dialog nodes or choices are available:

### Variable Conditions
```json
{
  "type": "variable",
  "name": "met_elder", 
  "operator": "equals",
  "value": true
}
```

### Level Conditions
```json
{
  "type": "level",
  "operator": "greater_than",
  "value": 10
}
```

### Item Conditions  
```json
{
  "type": "has_item",
  "item_id": "magic_sword",
  "quantity": 1
}
```

### Quest Conditions
```json
{
  "type": "quest_status",
  "quest_id": "rescue_princess",
  "status": "completed"
}
```

### Job/Class Conditions
```json
{
  "type": "job",
  "job_type": "warrior"
}
```

## Action System

Actions are executed when dialog nodes are reached:

### Set Variable
```json
{
  "type": "set_variable",
  "name": "elder_trust",
  "value": 5
}
```

### Add Experience
```json
{
  "type": "add_experience", 
  "amount": 100
}
```

### Give Item
```json
{
  "type": "give_item",
  "item_id": "health_potion",
  "quantity": 3
}
```

### Start Quest
```json
{
  "type": "start_quest",
  "quest_id": "find_missing_child"
}
```

### Play Sound
```json
{
  "type": "play_sound",
  "sound_id": "quest_complete"
}
```

## Choice Structure

Player choices define branching dialog paths:

```json
{
  "text": "I accept your quest!",
  "target": "quest_accepted",
  "conditions": [
    {"type": "level", "operator": "greater_than", "value": 3}
  ],
  "actions": [
    {"type": "set_variable", "name": "accepted_quest", "value": true}
  ]
}
```

### Choice Fields
- **`text`**: The choice text displayed to the player
- **`target`**: The dialog node to go to when selected
- **`conditions`**: Optional conditions to show this choice
- **`actions`**: Optional actions to execute when choice is selected

## Advanced Features

### Text Variables

Use variables in dialog text with `${variable_name}` syntax:

```json
{
  "text": "Greetings, ${player_name}! Your reputation as a ${player_class} precedes you."
}
```

### Conditional Text Blocks

Include conditional text within dialog:

```json
{
  "text": "Welcome back, traveler. {{if elder_trust > 2}}I'm glad to see a trusted friend return.{{else}}I hope your journey has been profitable.{{endif}}"
}
```

### Dialog Flags

Set persistent flags that affect future conversations:

```json
{
  "actions": [
    {"type": "set_flag", "flag": "completed_tutorial", "value": true}
  ]
}
```

## Example Complete Dialog

Here's a complete example dialog file:

```json
{
  "dialog_id": "merchant_greeting",
  "title": "Merchant Encounter",
  "description": "Meeting a traveling merchant",
  "variables": {
    "merchant_discount": 0,
    "bought_items": false
  },
  "nodes": {
    "start": {
      "speaker": "merchant",
      "text": "Well met, traveler! Care to browse my wares?",
      "choices": [
        {
          "text": "What do you have for sale?",
          "target": "show_shop"
        },
        {
          "text": "Any news from the road?",
          "target": "road_news"
        },
        {
          "text": "Not today, thanks.",
          "target": "polite_decline"
        }
      ]
    },
    "show_shop": {
      "speaker": "merchant",
      "text": "I have the finest goods from across the realm! Special prices for heroes like yourself.",
      "actions": [
        {"type": "open_shop", "shop_id": "traveling_merchant"}
      ],
      "end": true
    },
    "polite_decline": {
      "speaker": "merchant", 
      "text": "No worries, friend. Safe travels!",
      "end": true
    }
  }
}
```

## Integration with Game Systems

### Variable Persistence
- Dialog variables are saved with game state
- Variables can be accessed by other game systems
- Global flags affect multiple dialog trees

### Event Integration
- Dialog actions can trigger game events
- Game events can modify dialog availability
- Real-time game state affects dialog conditions

### Asset Management
- Character portraits loaded from `assets/portraits/`
- Sound effects loaded from `assets/audio/dialog/`
- Dialog scripts cached for performance

This format provides a powerful and flexible system for creating rich, branching narratives that integrate seamlessly with the game's RPG systems.