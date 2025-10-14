# Mermaid Diagrams

This directory contains Mermaid diagram files (.mmd) for visualizing various aspects of the MyRPG game architecture and flows.

## Files

- `tactical_combat_flow.md` - Complete tactical combat system documentation with Mermaid flowchart, ASCII diagrams, phase transitions, and UI flow explanations

This directory contains both pure Mermaid diagram files (.mmd) and comprehensive documentation files (.md) that include Mermaid diagrams alongside detailed explanations.

## How to Use

### Option 1: Mermaid Live Editor (Recommended)
1. Go to [mermaid.live](https://mermaid.live/)
2. Copy the content of any `.mmd` file
3. Paste it in the editor
4. Export as PNG, SVG, or PDF

### Option 2: GitHub Integration
- GitHub natively supports Mermaid diagrams
- Include in markdown files using:
```markdown
```mermaid
[paste .mmd file content here]
```
```

### Option 3: VS Code Extensions
- Install "Mermaid Markdown Syntax Highlighting" extension
- Install "Mermaid Preview" extension for live preview

### Option 4: Command Line Tools
```bash
# Install mermaid-cli
npm install -g @mermaid-js/mermaid-cli

# Generate PNG from .mmd file
mmdc -i tactical_combat_flow.mmd -o tactical_combat_flow.png

# Generate SVG
mmdc -i tactical_combat_flow.mmd -o tactical_combat_flow.svg
```

## Diagram Types Available

- **Flowcharts** - Process flows and decision trees
- **Sequence Diagrams** - Interaction between components over time
- **State Diagrams** - State machine transitions
- **Class Diagrams** - Object relationships and hierarchies

## Contributing

When adding new diagrams:
1. Use descriptive filenames (e.g., `ui_state_machine.mmd`, `ecs_component_flow.mmd`)
2. Include comments in the diagram for clarity
3. Update this README with the new file description
4. Test the diagram in Mermaid Live Editor before committing