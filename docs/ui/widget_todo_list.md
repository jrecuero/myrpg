# Widget Development Todo List

This document outlines recommended widgets for the RPG game, organized by implementation priority and complexity. Each widget builds on the existing popup widget infrastructure with constants-based customization.

## 🚀 High Priority Widgets (Essential for RPG gameplay)

### 1. ✅ [COMPLETED] Popup Selection Widget
- **Status**: ✅ COMPLETED
- **Features**: Scrollable list, arrow navigation, Enter/Escape handling, callback system
- **Integration**: UIManager, input blocking, P key trigger
- **Files**: `internal/ui/popup_selection_widget.go`

### 2. ✅ [COMPLETED] Popup Info Widget  
- **Status**: ✅ COMPLETED
- **Features**: Multi-line text display, scrollable content, ESC to close
- **Integration**: UIManager, input blocking, I key trigger
- **Files**: `internal/ui/popup_info_widget.go`

### 3. 📋 Inventory Widget
- **Status**: ⏳ TODO - High Priority
- **Features**: 
  - Grid-based item display with configurable slot sizes
  - Drag & drop item management
  - Item tooltips with stats and descriptions
  - Sorting options (type, value, name, rarity)
  - Search/filter functionality
  - Item stacking for consumables
- **Integration**: 
  - Triggered by 'I' key in exploration mode
  - Links with Equipment Widget for item equipping
  - Uses popup widget base architecture
  - Item data structure integration
- **Implementation Notes**:
  - Use constants for grid sizing, slot dimensions, tooltip styling
  - Support different inventory sizes based on character progression
  - Include visual indicators for item quality/rarity
- **Estimated Complexity**: Medium-High
- **Dependencies**: Item system, Equipment system

### 4. ✅ [COMPLETED] Character Stats Widget
- **Status**: ✅ COMPLETED
- **Features**: 
  - ✅ Detailed stat breakdown organized by categories (Overview, Core, Combat, Tactical)
  - ✅ Level progression display with XP bars
  - ✅ Job-specific stat organization and display
  - ✅ Visual progress bars for HP, MP, and XP
  - ✅ Tabbed interface with arrow key navigation
  - ✅ Color-coded sections for different stat categories
- **Integration**: 
  - ✅ Triggered by 'C' key in exploration mode  
  - ✅ Uses existing RPGStatsComponent data
  - ✅ UIManager integration with input blocking
  - ✅ JobType enum integration for job display
  - ✅ Engine integration with showCharacterStats() method
- **Files**: `internal/ui/character_stats_widget.go`, `test/character_stats_test/main.go`
- **Documentation**: `docs/ui/character_stats_constants.md`
- **Testing**: `make test-character-stats` - Interactive test program
- **Completed**: All features implemented with comprehensive constants system

### 5. ✅ [COMPLETED] Equipment Widget
- **Status**: ✅ COMPLETED - High Priority
- **Features**:
  - ✅ Paperdoll-style equipment slots (Weapon, Armor, Accessories)
  - ✅ Visual equipment preview with character model
  - ✅ Basic stat comparison framework in place
  - ✅ Equipment set bonuses display structure ready
  - ✅ Quick equipment swapping (ENTER key functionality)
- **Integration**:
  - ✅ Triggered by 'E' key in exploration mode
  - ✅ Mock equipment system for testing equip/unequip
  - ✅ Job restriction and level requirement validation
  - ✅ Equipment slot restrictions based on JobType
  - ✅ UIManager integration with input blocking
  - ✅ Engine integration with showEquipment() method
- **Implementation Notes**:
  - ✅ Use slot-based layout with visual equipment representation
  - ✅ Color-code stat changes (green=better, red=worse) (colors defined)
  - ✅ Support equipment requirements (level, stats, job)
  - ✅ ESC key handling and navigation system
  - ✅ Contextual help text showing current action
- **Files**: `internal/ui/equipment_widget.go`, `internal/ecs/components/equipment.go`, `test/equipment_test/main.go`
- **Documentation**: `docs/ui/equipment_constants.md`
- **Testing**: `make test-equipment` - Interactive test program with equip/unequip
- **Engine Integration**: ✅ Complete - 'E' key trigger, UIManager methods, Entity component access
- **Completed**: Basic equip/unequip functionality, job restrictions, mock equipment system
- **Future Enhancements**: Advanced stat comparison, inventory integration, equipment sets
- **Estimated Complexity**: Medium-High
- **Dependencies**: ✅ Complete for basic functionality

### 6. ✅ [COMPLETED] Dialog Widget
- **Status**: ✅ COMPLETED - Fully integrated with main game
- **Features**:
  - ✅ External dialog script system (JSON format)
  - ✅ NPC conversation display with portraits
  - ✅ Branching dialog trees with conditions
  - ✅ Multiple choice selection trees
  - ✅ Game event integration and variable system
  - ✅ Typewriter text effect for immersive display
  - ⏳ Text formatting (bold, italic, colors) - Future enhancement
  - ⏳ Dialog history/log - Future enhancement
- **Integration**:
  - ✅ UIManager integration with input blocking
  - ✅ Triggered by D key in main game engine
  - ✅ Game state conditions and variable tracking
  - ⏳ Save/load dialog state for complex conversations - Future enhancement
- **Implementation Notes**:
  - External JSON script files for dialog content
  - Conditional branching based on game events
  - Portrait positioning and character display
  - Follows established popup widget patterns
- **Files**: `internal/ui/dialog_widget.go`, `assets/dialogs/*.json`, `test/dialog_test/main.go`
- **Documentation**: `docs/ui/dialog_constants.md`, `docs/ui/dialog_script_format.md`, `docs/ui/dialog_integration.md`
- **Testing**: `make test-dialog` - Interactive dialog test with sample scripts
- **Engine Integration**: `internal/engine/engine.go` - D key trigger, showDialog() method
- **Completed**: Full dialog system with external JSON scripts, branching conversations, variable system
- **Future Enhancements**: Dynamic NPC selection, save/load state, quest integration

## 📈 Medium Priority Widgets (Enhances experience)

### 7. 🎯 Skills/Abilities Widget
- **Status**: ⏳ TODO - Medium Priority
- **Features**:
  - Job-specific skill trees with branching paths
  - Skill point allocation and preview
  - Ability tooltips with damage/effects
  - Cooldown and resource cost tracking
  - Skill prerequisites visualization
- **Integration**:
  - Triggered by 'K' key (Skills)
  - Uses JobType system for skill tree selection
  - Links with Character Stats for skill point availability
  - Combat integration for ability usage
- **Implementation Notes**:
  - Tree-like layout with connecting lines
  - Interactive nodes for skill point allocation
  - Preview mode for planning builds
- **Estimated Complexity**: High
- **Dependencies**: Skill system, Job progression system

### 8. 📋 Quest Journal Widget
- **Status**: ⏳ TODO - Medium Priority
- **Features**:
  - Active quests with objective tracking
  - Completed quests archive
  - Quest categories (Main, Side, Daily)
  - Search and filtering options
  - Quest location hints/waypoints
- **Integration**:
  - Triggered by 'J' key (Journal)
  - Links with Dialog Widget for quest acceptance
  - Map integration for quest locations
  - Progress tracking and notifications
- **Implementation Notes**:
  - Tabbed interface for quest categories
  - Progress bars for multi-step quests
  - Rich text for quest descriptions
- **Estimated Complexity**: Medium
- **Dependencies**: Quest system, Map/Location system

### 9. 🎁 Loot/Reward Widget
- **Status**: ⏳ TODO - Medium Priority
- **Features**:
  - Item rewards display with rarity highlighting
  - Experience and gold gain notifications
  - Loot distribution for party members
  - Auto-pickup options with filters
  - Loot history and statistics
- **Integration**:
  - Appears after combat victories
  - Triggered by treasure chest interactions
  - Links with Inventory Widget for item collection
  - Party distribution mechanics
- **Implementation Notes**:
  - Animated item reveals for excitement
  - Rarity color coding and effects
  - Batch collection options
- **Estimated Complexity**: Medium
- **Dependencies**: Loot system, Party management

## 🔧 Lower Priority Widgets (Nice to have)

### 10. 📈 Experience/Progress Widget
- **Status**: ⏳ TODO - Lower Priority
- **Features**:
  - Detailed XP breakdown by source (combat, quests, exploration)
  - Level progress visualization with milestones
  - Recent gains log with timestamps
  - Next level stat preview
  - Progress statistics and achievements
- **Integration**:
  - Overlay display during XP gains
  - Part of Character Stats Widget as sub-panel
  - Uses existing Experience/ExpToNext system
  - Achievement system integration
- **Implementation Notes**:
  - Smooth progress bar animations
  - XP source categorization and tracking
  - Level-up celebration effects
- **Estimated Complexity**: Low-Medium
- **Dependencies**: Existing XP system, Achievement system

### 11. 🗺️ Minimap Widget
- **Status**: ⏳ TODO - Lower Priority
- **Features**:
  - World overview with fog of war
  - Unit positions (allies, enemies, NPCs)
  - Interactive map navigation
  - Zoom levels and detail modes
  - Waypoint and marker system
- **Integration**:
  - Toggleable overlay (M key)
  - Works with tactical grid system
  - Links with Quest Journal for objective markers
  - Exploration mode navigation aid
- **Implementation Notes**:
  - Efficient rendering for large maps
  - Real-time unit tracking
  - Customizable visibility options
- **Estimated Complexity**: High
- **Dependencies**: Map system, World generation

### 12. ⚙️ Settings/Options Widget
- **Status**: ⏳ TODO - Lower Priority
- **Features**:
  - Game options organized by categories
  - Customizable key bindings with conflict detection
  - Audio and video settings
  - Gameplay preferences (auto-save, difficulty)
  - Import/export settings profiles
- **Integration**:
  - Accessible from main menu and ESC key
  - Settings persistence system
  - Real-time preview for visual changes
  - Input validation and conflict resolution
- **Implementation Notes**:
  - Tabbed interface for setting categories
  - Live preview for changes
  - Reset to defaults functionality
- **Estimated Complexity**: Medium
- **Dependencies**: Settings persistence, Configuration system

## 🏗️ Implementation Guidelines

### Architecture Patterns
- **Base Class**: All widgets should extend popup widget architecture
- **Constants**: Use comprehensive constants for all visual properties
- **Input Blocking**: Integrate with existing UIManager input blocking system
- **Callbacks**: Implement callback system for widget interactions
- **Testing**: Create test programs for each widget (following popup widget pattern)

### Code Organization
```
internal/ui/
├── popup_base_widget.go          # Common popup functionality
├── inventory_widget.go            # Inventory management
├── character_stats_widget.go      # Character information
├── equipment_widget.go            # Equipment management
├── dialog_widget.go               # NPC conversations
├── skills_widget.go               # Skill trees and abilities
├── quest_journal_widget.go        # Quest tracking
├── loot_widget.go                 # Reward displays
├── progress_widget.go             # Experience tracking
├── minimap_widget.go              # World navigation
├── settings_widget.go             # Game options
└── ui_manager.go                  # Central widget coordination
```

### Constants Organization
```
docs/ui/
├── popup_constants.md             # ✅ COMPLETED - Popup widget constants
├── character_stats_constants.md   # ✅ COMPLETED - Character stats constants
├── inventory_constants.md          # Grid sizes, slot dimensions
├── equipment_constants.md          # Slot positions, preview areas
└── widget_constants_master.md      # Master reference for all widgets
```

### Testing Strategy
```
test/
├── inventory_test/                # Inventory widget functionality
├── character_stats_test/          # Stats display and navigation
├── equipment_test/                # Equipment preview and swapping
└── widget_integration_test/       # Multi-widget interactions
```

## 📋 Implementation Checklist Template

For each widget implementation:
- [ ] Create widget file with comprehensive constants
- [ ] Implement core functionality with popup base
- [ ] Add UIManager integration with input blocking
- [ ] Create test program for standalone testing
- [ ] Add Makefile target for easy testing
- [ ] Update documentation with constants reference
- [ ] Integrate with existing game systems
- [ ] Add keyboard shortcuts and hotkeys
- [ ] Test input conflicts and resolution
- [ ] Verify visual layout and responsiveness

## 🎯 Next Steps

1. **Choose Implementation Order**: Start with Inventory Widget as highest impact
2. **Set Up Base Architecture**: Create common popup base class if needed
3. **Define Data Structures**: Ensure Item, Equipment systems are ready
4. **Plan Integration Points**: Map out UIManager coordination
5. **Create Development Branches**: One widget per branch for clean development

## 📝 Notes

- All widgets should follow the established popup widget patterns
- Prioritize user experience and intuitive navigation
- Maintain consistent visual style with existing UI
- Ensure proper input handling and conflict resolution
- Plan for future extensibility and customization options
- Consider accessibility and usability in all designs

---

This todo list serves as a comprehensive roadmap for widget development. Each widget should be implemented incrementally, tested thoroughly, and integrated smoothly with the existing game systems.