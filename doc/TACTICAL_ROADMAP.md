# Phase 2: Movement Range System & Tactical Positioning

**Target**: Implement movement ranges, pathfinding, and basic tactical positioning

## 🎯 **Objectives**

### **Core Features**
- Calculate movement ranges based on character stats
- Visual highlighting of valid movement tiles
- Pathfinding with movement cost calculation
- Click-to-move tile selection
- Movement preview and confirmation

### **Technical Components**

#### **Movement Range Calculation**
```go
// Enhanced TacticalStatsComponent
type TacticalStatsComponent struct {
    *RPGStatsComponent       // Inherit existing stats
    MoveRange    int         // Tiles per turn (Warrior: 3, Mage: 2, Rogue: 4)
    MoveCost     int         // Cost per tile (usually 1)
    JumpHeight   int         // Max height difference (future: 0 for 2D)
}

// Movement calculation
func (tc *TacticalCombat) CalculateMovementRange(unit *ecs.Entity) []GridPos {
    stats := unit.TacticalStats()
    currentPos := unit.GridPosition()
    return tc.Grid.CalculateMovementRange(currentPos, stats.MoveRange)
}
```

#### **Pathfinding System**
```go
// A* pathfinding for optimal movement routes
type PathFinder struct {
    Grid *Grid
}

func (pf *PathFinder) FindPath(start, end GridPos, moveRange int) ([]GridPos, bool) {
    // A* algorithm implementation
    // Returns path and whether destination is reachable
}

func (pf *PathFinder) CalculateMoveCost(from, to GridPos) int {
    // Different terrain types have different costs
    // Basic: 1 per tile, water: 2, walls: impassable
}
```

#### **Grid Position Component**
```go
// New component to track tactical positions
type GridPositionComponent struct {
    Position GridPos    // Current grid coordinates
    Facing   Direction  // Unit facing (N, S, E, W)
}

// Conversion between world and grid coordinates
func (gpc *GridPositionComponent) SyncWithTransform(transform *Transform, grid *Grid) {
    // Keep world position and grid position synchronized
}
```

### **User Experience**

#### **Movement Flow**
1. **Unit Selection**: Click on unit to select (yellow highlight)
2. **Range Display**: Blue tiles show valid movement range
3. **Path Preview**: Green tiles show movement path on hover
4. **Confirmation**: Click destination to confirm movement
5. **Animation**: Unit slides smoothly between tiles

#### **Visual Feedback**
- **Blue Tiles**: Valid movement destinations
- **Green Path**: Hover preview of movement route
- **Yellow Selection**: Currently selected unit
- **Red X**: Invalid/blocked tiles
- **Movement Counter**: Display remaining movement points

### **Controls**
- **Mouse Click**: Select unit or destination tile
- **Hover**: Preview movement path
- **Right Click**: Cancel selection
- **Space**: Confirm movement
- **C**: Cancel current action

## 📋 **Implementation Steps**

### **Week 1: Foundation**
1. Create GridPositionComponent
2. Add TacticalStatsComponent 
3. Implement basic movement range calculation
4. Add mouse input handling for tile selection

### **Week 2: Movement System**
1. Implement A* pathfinding algorithm
2. Add movement cost calculation
3. Create path preview system
4. Add tile selection and highlighting

### **Week 3: Animation & Polish**
1. Smooth movement animations between tiles
2. Movement confirmation system
3. Visual feedback improvements
4. Mouse cursor changes and tooltips

### **Week 4: Integration**
1. Integrate with existing character system
2. Add movement to tactical combat flow
3. Testing and bug fixes
4. Performance optimization

---

# Phase 3: Turn-Based Combat System

**Target**: Replace collision-based combat with full turn-based tactical combat

## 🎯 **Objectives**

### **Core Features**
- Initiative-based turn order system
- Action phases (Move → Action → End Turn)
- Turn queue visualization
- Action point system
- Combat state management

### **Technical Components**

#### **Turn Management System**
```go
type TurnManager struct {
    TurnQueue     []*CombatUnit    // Ordered by initiative
    CurrentTurn   int              // Index in turn queue
    TurnCounter   int              // Total turns elapsed
    Phase         TurnPhase        // Current phase of turn
}

type TurnPhase int
const (
    TurnPhaseMove TurnPhase = iota    // Movement phase
    TurnPhaseAction                   // Action selection phase
    TurnPhaseConfirm                  // Confirmation phase
    TurnPhaseExecute                  // Action execution
    TurnPhaseEnd                      // End turn cleanup
)

type CombatUnit struct {
    Entity      *ecs.Entity
    Initiative  int              // Turn order priority
    ActionPoints int             // Actions remaining this turn
    HasMoved     bool            // Movement used this turn
    HasActed     bool            // Action used this turn
    GridPos      GridPos         // Current position
}
```

#### **Initiative System**
```go
func CalculateInitiative(stats *RPGStatsComponent) int {
    base := stats.Speed
    jobBonus := getJobInitiativeBonus(stats.Job)
    randomFactor := rand.Intn(10) // 0-9 random
    return base + jobBonus + randomFactor
}

func getJobInitiativeBonus(job JobType) int {
    switch job {
    case JobRogue:   return 3  // Fast
    case JobWarrior: return 1  // Average  
    case JobMage:    return 0  // Slow
    default:         return 0
    }
}
```

#### **Action Point System**
```go
type ActionCost struct {
    Move   int  // Cost to move
    Attack int  // Cost to attack
    Skill  int  // Cost to use skill
    Item   int  // Cost to use item
    Wait   int  // Cost to wait/end turn
}

var DefaultActionCosts = ActionCost{
    Move:   1,
    Attack: 1, 
    Skill:  2,
    Item:   1,
    Wait:   0,
}
```

### **Combat Flow**

#### **Turn Sequence**
1. **Initiative Roll**: Calculate turn order at combat start
2. **Turn Start**: Select next unit in queue
3. **Movement Phase**: Player selects movement (optional)
4. **Action Phase**: Player selects action (attack, skill, item, wait)
5. **Confirmation**: Preview action effects and confirm
6. **Execution**: Execute action with animations
7. **Turn End**: Apply effects, check victory conditions
8. **Next Turn**: Move to next unit in queue

#### **Victory Conditions**
- **Defeat All Enemies**: Classic victory condition
- **Survive X Turns**: Defensive scenarios
- **Reach Target**: Movement-based objectives
- **Protect VIP**: Escort missions

### **User Interface**

#### **Turn Order Display**
```go
type TurnOrderUI struct {
    Units        []*CombatUnit    // All units in combat
    CurrentIndex int              // Highlight current turn
    Position     image.Point      // UI panel position
}

// Display portraits in turn order with initiative values
func (tui *TurnOrderUI) Draw(screen *ebiten.Image) {
    // Portrait + initiative number for each unit
    // Highlight current turn with border
    // Show action points remaining
}
```

#### **Action Menu System**
```go
type ActionMenu struct {
    Actions     []ActionOption   // Available actions
    Selected    int             // Currently selected action
    Visible     bool            // Menu visibility
}

type ActionOption struct {
    Name        string
    Icon        *ebiten.Image
    Cost        int             // Action point cost
    Available   bool            // Can be used
    Hotkey      ebiten.Key      // Keyboard shortcut
}
```

## 📋 **Implementation Steps**

### **Week 1: Turn System Foundation**
1. Create TurnManager and CombatUnit structures
2. Implement initiative calculation
3. Build turn queue management
4. Add basic turn progression

### **Week 2: Action System**
1. Create action point system
2. Implement action costs and validation
3. Add action menu UI
4. Create action confirmation system

### **Week 3: Combat Integration**
1. Replace collision battles with tactical combat
2. Add combat start/end conditions
3. Implement victory/defeat detection
4. Add turn order UI display

### **Week 4: Polish & Testing**
1. Animation timing and visual feedback
2. AI behavior for enemy turns
3. Balance testing and adjustment
4. Performance optimization

---

# Phase 4: Action System & Combat Mechanics

**Target**: Implement comprehensive action system with attacks, skills, and items

## 🎯 **Objectives**

### **Core Features**
- Attack range and targeting system
- Skill/magic system with area effects
- Item usage in combat
- Damage calculation with tactical modifiers
- Status effects and buffs/debuffs

### **Technical Components**

#### **Targeting System**
```go
type TargetingSystem struct {
    AttackRange   map[JobType]int       // Attack ranges by job
    SkillRanges   map[string]AreaEffect // Skill targeting areas
    LineOfSight   bool                  // LOS calculations
}

type AreaEffect struct {
    Shape    EffectShape    // Single, Line, Cross, Square, Circle
    Size     int           // Area size
    Range    int           // Cast range
    Friendly bool          // Affects allies
    Enemy    bool          // Affects enemies
}

type EffectShape int
const (
    EffectSingle EffectShape = iota
    EffectLine                // Linear effect
    EffectCross               // + shape
    EffectSquare              // Square area
    EffectCircle              // Circular area
)
```

#### **Enhanced Combat System**
```go
type TacticalDamageCalculator struct {
    BaseDamage    int
    Modifiers     []DamageModifier
}

type DamageModifier struct {
    Type      ModifierType
    Value     float64
    Source    string
}

type ModifierType int
const (
    ModifierHeight     ModifierType = iota  // Height advantage
    ModifierFacing                          // Back/side attacks
    ModifierDistance                        // Range penalties
    ModifierTerrain                         // Terrain effects
    ModifierStatus                          // Status effect bonuses
)

func (tdc *TacticalDamageCalculator) Calculate(attacker, defender *CombatUnit) int {
    base := attacker.Entity.RPGStats().Attack - defender.Entity.RPGStats().Defense
    
    // Apply tactical modifiers
    for _, modifier := range tdc.Modifiers {
        base = int(float64(base) * modifier.Value)
    }
    
    return max(1, base) // Always at least 1 damage
}
```

#### **Skill System**
```go
type SkillComponent struct {
    Skills       map[string]*Skill    // Available skills
    MPCost       map[string]int       // MP costs
    Cooldowns    map[string]int       // Turn cooldowns
    UsesPerBattle map[string]int      // Limited uses
}

type Skill struct {
    Name         string
    Description  string
    Range        int
    Area         AreaEffect
    Effects      []SkillEffect
    Animation    string
    MPCost       int
    Cooldown     int
    UsesPerBattle int
}

type SkillEffect struct {
    Type     EffectType
    Value    int
    Duration int        // Turns for status effects
    Target   TargetType // Self, Enemy, Ally, All
}

type EffectType int
const (
    EffectDamage EffectType = iota
    EffectHeal
    EffectBuff
    EffectDebuff
    EffectStatus
)
```

### **Job-Specific Skills**

#### **Warrior Skills**
- **Charge**: Move + attack in one action
- **Guard**: Reduce damage, protect adjacent allies
- **Taunt**: Force enemies to target this unit
- **Weapon Break**: Reduce enemy attack power

#### **Mage Skills**
- **Fireball**: High damage, small area
- **Heal**: Restore HP to ally
- **Ice**: Damage + slow effect
- **Barrier**: Magical damage reduction

#### **Rogue Skills**
- **Backstab**: High damage from behind
- **Stealth**: Become untargetable for 1 turn
- **Poison**: Damage over time effect
- **Lockpick**: Open chests/doors in tactical maps

### **Item System**
```go
type TacticalItem struct {
    Name        string
    Type        ItemType
    Effect      ItemEffect
    Range       int
    Area        AreaEffect
    Consumable  bool
    Quantity    int
}

type ItemType int
const (
    ItemPotion ItemType = iota
    ItemWeapon
    ItemArmor
    ItemTool
    ItemKey
)

type ItemEffect struct {
    HPRestore   int
    MPRestore   int
    StatusCure  []StatusType
    StatBoost   map[string]int
    Duration    int
}
```

## 📋 **Implementation Steps**

### **Week 1: Targeting Foundation**
1. Implement attack range calculation
2. Add target selection system
3. Create line-of-sight validation
4. Build area effect targeting

### **Week 2: Skill System**
1. Create skill component and data structures
2. Implement MP cost and cooldown system
3. Add skill effects and status system
4. Create skill selection UI

### **Week 3: Enhanced Combat**
1. Add tactical damage modifiers
2. Implement facing and positioning bonuses
3. Create status effect system
4. Add combat animations for skills

### **Week 4: Item Integration**
1. Add item usage in combat
2. Create inventory management for tactical mode
3. Implement consumable items
4. Add equipment effects in tactical combat

---

# Phase 5: Advanced Tactical Features

**Target**: Add advanced FFT-style mechanics and polish

## 🎯 **Objectives**

### **Core Features**
- Terrain effects and environmental interactions
- Advanced AI for enemy tactics
- Job system integration with tactical abilities
- Save/Load tactical combat states
- Multiplayer foundation

### **Advanced Mechanics**

#### **Environmental System**
```go
type TerrainEffect struct {
    Type         TerrainType
    MoveCost     int          // Movement cost modifier
    DefenseBonus int          // Defense bonus when standing on
    StatusEffect StatusType   // Ongoing effect (poison swamp, etc.)
    Special      []string     // Special interactions
}

type TerrainType int
const (
    TerrainGrass TerrainType = iota
    TerrainWater                    // Slow movement
    TerrainSwamp                    // Poison + slow
    TerrainHill                     // Defense bonus
    TerrainWall                     // Impassable
    TerrainTrap                     // Hidden damage
    TerrainHealing                  // HP regeneration
    TerrainMagic                    // MP regeneration
)
```

#### **AI Tactical System**
```go
type TacticalAI struct {
    Difficulty   AILevel
    Personality  AIPersonality
    Objectives   []AIObjective
}

type AILevel int
const (
    AIBeginner AILevel = iota
    AIIntermediate
    AIAdvanced
    AIExpert
)

type AIPersonality int
const (
    AIAggressive AIPersonality = iota  // Focus on attacks
    AIDefensive                        // Focus on positioning
    AISupport                          // Focus on buffs/heals
    AIBalanced                         // Mixed strategy
)

type AIObjective struct {
    Type     ObjectiveType
    Target   *ecs.Entity
    Priority int
    Condition string
}
```

#### **Advanced Job System**
```go
type JobAdvancement struct {
    CurrentJob    JobType
    Level         int
    Experience    int
    UnlockedSkills []string
    StatGrowth    map[string]float64
}

type JobUnlock struct {
    RequiredJobs  []JobType
    RequiredLevel int
    SpecialReq    string
}

// Advanced jobs unlocked through gameplay
const (
    JobPaladin   JobType = iota + 10  // Warrior + Mage hybrid
    JobNinja                          // Advanced Rogue
    JobSorcerer                       // Advanced Mage
    JobBerserker                      // High-risk Warrior
)
```

### **Quality of Life Features**

#### **Tactical Map Editor**
```go
type MapEditor struct {
    Grid         *Grid
    TileSet      []TileType
    UnitPlacer   *UnitPlacer
    ObjectPlacer *ObjectPlacer
    SaveLoad     *MapSaveLoad
}

// Allow custom tactical maps
func (me *MapEditor) CreateCustomMap(width, height int) *TacticalMap {
    // Map creation interface
}
```

#### **Replay System**
```go
type TacticalReplay struct {
    Actions    []CombatAction
    GameState  []GameSnapshot
    Metadata   ReplayMetadata
}

type CombatAction struct {
    Turn      int
    Unit      string
    Action    string
    Target    GridPos
    Result    ActionResult
    Timestamp time.Time
}
```

## 📋 **Implementation Timeline**

### **Month 1: Environmental System**
- Terrain types and effects
- Environmental interactions
- Weather and dynamic conditions
- Map hazards and traps

### **Month 2: AI Enhancement**
- Tactical AI behaviors
- Difficulty scaling
- AI personality systems
- Performance optimization

### **Month 3: Advanced Features**
- Job advancement system
- Special abilities and ultimates
- Equipment and item crafting
- Achievement system

### **Month 4: Polish & Extras**
- Map editor and custom scenarios
- Replay system
- Multiplayer foundation
- Performance optimization and testing

---

# Phase 6: Content & Polish

**Target**: Add content, balance, and final polish for release

## 🎯 **Objectives**

### **Content Creation**
- Campaign with tactical scenarios
- Character progression and story
- Multiple victory conditions
- Difficulty options and accessibility

### **Balance & Testing**
- Comprehensive playtesting
- Statistical balance analysis
- Performance optimization
- Bug fixing and stability

### **Release Preparation**
- Documentation and tutorials
- Save system compatibility
- Platform compatibility
- Marketing materials

---

**This comprehensive roadmap gives you a complete path from your current excellent foundation to a full-featured FFT-style tactical RPG!** 🗡️⚔️✨