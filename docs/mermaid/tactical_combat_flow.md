# Tactical Combat Flow Diagram

## Mermaid Flowchart

```mermaid
flowchart TD
    A[Game Start] --> B[Enter Tactical Mode]
    B --> C[Initialize Combat]
    
    C --> D[CombatPhaseInitialization]
    D --> D1[Setup Teams & Initiative]
    D1 --> D2[Deploy Units on Grid]
    D2 --> D3[Calculate Turn Order]
    D3 --> E[CombatPhaseTeamTurn]
    
    E --> E1[Set Active Team]
    E1 --> E2[Restore Team AP]
    E2 --> E3[Select Active Unit]
    E3 --> E4{Player Input}
    
    E4 -->|Arrow Keys| F[Move Unit]
    E4 -->|A Key| G[Attack Action]
    E4 -->|E Key| H[End Turn]
    E4 -->|TAB Key| I[Switch Unit]
    
    F --> F1[Validate Movement]
    F1 --> F2{Valid Move?}
    F2 -->|Yes| F3[Update Position]
    F2 -->|No| E4
    F3 --> E4
    
    G --> G1[CombatPhaseActionExecution]
    G1 --> G2[Validate Attack]
    G2 --> G3{Valid Target?}
    G3 -->|Yes| G4[Calculate Damage]
    G3 -->|No| E4
    G4 --> G5[Apply Damage]
    G5 --> G6[Update UI Messages]
    G6 --> E4
    
    I --> I1{More Units Available?}
    I1 -->|Yes| I2[Switch to Next Unit]
    I1 -->|No| E4
    I2 --> E4
    
    H --> J[CombatPhaseEndTurn]
    J --> J1[Mark Team Complete]
    J1 --> J2{All Teams Complete?}
    J2 -->|No| J3[Next Team Turn]
    J2 -->|Yes| K[CombatPhaseVictoryCheck]
    J3 --> E
    
    K --> K1[Check Win Conditions]
    K1 --> K2{Combat Over?}
    K2 -->|Player Victory| L1[CombatPhaseEnded - Player Wins]
    K2 -->|Enemy Victory| L2[CombatPhaseEnded - Enemy Wins]
    K2 -->|Continue| K3[Increment Round]
    K3 --> E
    
    L1 --> M[Return to Exploration]
    L2 --> N[Game Over]
    
    style D fill:#e1f5fe
    style E fill:#f3e5f5
    style G1 fill:#fff3e0
    style J fill:#e8f5e8
    style K fill:#fce4ec
    style L1 fill:#e0f2f1
    style L2 fill:#ffebee
```

## ASCII Flow Diagram

```
┌─────────────────┐
│   Game Start    │
└─────────┬───────┘
          │
┌─────────▼───────┐
│ Enter Tactical  │
│     Mode        │
└─────────┬───────┘
          │
┌─────────▼───────────────────────────────────┐
│         INITIALIZATION PHASE                │
│ • Setup teams & initiative order           │
│ • Deploy units on grid                     │
│ • Calculate turn sequence                  │
└─────────┬───────────────────────────────────┘
          │
┌─────────▼───────────────────────────────────┐
│           TEAM TURN PHASE                   │
│ • Set active team & restore AP             │
│ • Select active unit                       │
│ • Handle player input:                     │
│   - Arrow Keys: Move unit                  │
│   - A Key: Attack action                   │
│   - E Key: End turn                        │
│   - TAB Key: Switch unit                   │
└─────────┬───────────────────────────────────┘
          │
          ├── Move ──┐
          │          │
          ├── Attack ─┼─► ACTION EXECUTION PHASE
          │          │   • Validate action
          │          │   • Calculate & apply effects
          │          │   • Update UI messages
          │          └─► Back to Team Turn
          │
          └── End Turn ──┐
                         │
┌────────────────────────▼─────────────────────┐
│            END TURN PHASE                    │
│ • Mark current team as complete             │
│ • Check if all teams finished round        │
│ • Switch to next team OR proceed to        │
│   victory check                             │
└─────────┬───────────────────────────────────┘
          │
┌─────────▼───────────────────────────────────┐
│          VICTORY CHECK PHASE                │
│ • Check win/lose conditions                │
│ • If combat continues: increment round     │
│ • If over: transition to ENDED phase       │
└─────────┬───────────────────────────────────┘
          │
     ┌────┴────┐
     │ Victory │ Defeat │ Continue
     ▼         ▼        ▼
┌─────────┐ ┌─────────┐ │
│ PLAYER  │ │  ENEMY  │ │
│ VICTORY │ │ VICTORY │ │
└─────────┘ └─────────┘ │
     │         │        │
     └─────────┼────────┘
               │ Loop back to
               │ TEAM TURN
               │ (new round)
```

## Phase Transitions & Callbacks

```
StateChangeCallback triggers on these transitions:

CombatPhaseInitialization ──► CombatPhaseTeamTurn
                                       │
                                       ▼
CombatPhaseActionExecution ◄── Player Input (Attack)
       │
       ▼
CombatPhaseTeamTurn ◄────────── Action Complete
       │
       ▼
CombatPhaseEndTurn ◄─────────── Player Press 'E'
       │
       ▼
CombatPhaseVictoryCheck
       │
       ├─► CombatPhaseEnded (if victory/defeat)
       │
       └─► CombatPhaseTeamTurn (if continue, new round)
```

## UI Updates Flow

```
Combat Phase Change
       │
       ▼ StateChangeCallback
Engine.updateMovementHighlighting()
       │
       ▼
Check if CombatPhaseTeamTurn
       │
       ▼
Get Active Player
       │
       ▼
Update Movement Range Display
       │
       ▼
UI Reflects Current Unit's Movement Options
```