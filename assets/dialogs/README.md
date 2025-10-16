# Dialog Trees Index

## Overview
Visual representations of all dialog trees in the game. These Mermaid diagrams help understand conversation flow, plan dialog expansions, and identify connection patterns.

## Available Dialog Trees

### 🏛️ Town Elder (Elder Aldric)
**File**: [town_elder_tree.md](town_elder_tree.md)
**JSON**: [town_elder.json](town_elder.json)
**Complexity**: High (19 nodes, multiple quest paths)

**Key Features**:
- Hero's journey path with trials
- Bandit quest line
- Village information and lore
- Trust system (0-4 levels)
- Multiple ending scenarios

**Main Paths**:
- 🗡️ Prophecy → Hero Claim → Trials → Quest
- 🏹 Work Request → Bandit Quest → Mission
- 📖 Village Info → Dangers → Investigation
- 🛡️ Protection → Help Offer → Guidance

---

### 🛒 Merchant (Trader Gareth)
**File**: [merchant_tree.md](merchant_tree.md)
**JSON**: [merchant.json](merchant.json)
**Complexity**: Medium (16 nodes, commerce focused)

**Key Features**:
- Shopping and trade interface
- Road news and rumors
- Bandit quest (alternate source)
- Reputation system
- Negotiation mechanics

**Main Paths**:
- 🛒 Wares → Shopping → Transaction
- 📰 Road News → Creatures → Quest Offer
- 🗡️ Bandit Trouble → Help Offer → Mission
- 💰 Negotiation → Payment → Accept Quest

---

## Dialog Tree Statistics

| Character | Nodes | End Nodes | Quest Paths | Variables | Complexity |
|-----------|-------|-----------|-------------|-----------|------------|
| Town Elder | 19 | 4 | 3 | 6 | High |
| Merchant | 16 | 3 | 1 | 4 | Medium |
| **Total** | **35** | **7** | **4** | **10** | - |

## Common Dialog Patterns

### 🔄 **Hub Pattern** (Town Elder)
```
start → multiple_topics → various_paths → multiple_endings
```
- **Pros**: Rich exploration, multiple replays
- **Cons**: Complex to maintain
- **Best For**: Major NPCs, story-critical characters

### 🎯 **Linear Pattern** (Merchant)
```
start → topic_selection → focused_path → conclusion
```
- **Pros**: Clear purpose, easy navigation
- **Cons**: Limited replayability
- **Best For**: Service NPCs, quest givers

### 🔄 **Circular Navigation**
Both dialogs use `other_topics` or similar nodes to return to main conversation branches, allowing players to explore multiple topics in one conversation.

## Color Coding Legend

| Color | Purpose | Examples |
|-------|---------|----------|
| 🟢 Green | Start nodes | `start` |
| 🔴 Red | End nodes | `end_friendly`, `polite_farewell` |
| 🟠 Orange | Quest nodes | `accept_bandit_quest`, `accept_challenge` |
| 🔵 Blue | Info nodes | `ask_village`, `road_news` |
| 🟣 Purple | Action nodes | `claim_hero`, `offer_help` |
| 🟤 Brown | Negotiation | `negotiate_payment`, `demand_more` |

## Planning New Dialog Trees

### 📋 **Recommended Process**
1. **Design Conversation Goals**: What should the player learn/achieve?
2. **Map Key Topics**: 3-5 main conversation branches
3. **Plan Quest Integration**: How do quests connect to conversation?
4. **Design Variable System**: What character development occurs?
5. **Create JSON Structure**: Implement dialog nodes
6. **Generate Mermaid Diagram**: Visualize and validate flow
7. **Test All Paths**: Ensure no broken connections

### 🎯 **Complexity Guidelines**
- **Simple (5-8 nodes)**: Service NPCs, simple merchants
- **Medium (9-16 nodes)**: Quest givers, story NPCs
- **Complex (17+ nodes)**: Major characters, hubs

### 🔗 **Connection Patterns**
- **Minimum**: Every node reachable from start
- **Recommended**: Multiple paths to important nodes
- **Advanced**: Circular navigation for topic exploration

## Future Expansions

### 🚀 **Planned Characters**
- **Guard Captain**: Security and training dialog
- **Blacksmith**: Equipment and crafting conversations  
- **Innkeeper**: Rumors and room rentals
- **Village Priest**: Blessings and spiritual guidance
- **Mysterious Stranger**: Advanced quest hooks

### 🎨 **Dialog Enhancements**
- **Conditional Branching**: Level/item-gated conversations
- **Dynamic Text**: Variable-driven dialog variations
- **Character Relationships**: Cross-character reputation effects
- **Quest Integration**: Dialog changes based on completed quests
- **Time-based Variations**: Different conversations at different times

## Maintenance

### 📅 **Regular Updates**
- ✅ **After JSON Changes**: Update corresponding tree diagram
- ✅ **New Characters**: Create tree diagram before implementation
- ✅ **Quest Integration**: Update affected character trees
- ✅ **Variable Changes**: Update documentation sections

### 🔍 **Validation Tools**
- **Mermaid Live Editor**: https://mermaid.live/
- **JSON Validator**: Ensure syntax correctness
- **Dialog Testing**: `make test-dialog` for runtime validation

### 📖 **Documentation Links**
- [Dialog Script Format](../docs/ui/dialog_script_format.md)
- [Character Management](../docs/ui/dialog_character_management.md)
- [Dialog Integration](../docs/ui/dialog_integration.md)
- [Tree Creation Guide](dialog_tree_guide.md)

This index provides a comprehensive overview of the dialog system's visual documentation and serves as a central reference for understanding character interactions!