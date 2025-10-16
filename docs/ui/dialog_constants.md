# Dialog Widget Constants

This document contains all the constants used in the Dialog Widget implementation. These values control layout, styling, colors, and behavior of the dialog system.

## Layout Dimensions

| Constant | Value | Description |
|----------|--------|-------------|
| `DialogWidgetWidth` | 700 | Widget width in pixels |
| `DialogWidgetHeight` | 400 | Widget height in pixels |
| `DialogWidgetBorderWidth` | 2 | Border thickness |
| `DialogWidgetShadowOffset` | 6 | Shadow offset distance |
| `DialogWidgetPadding` | 20 | Internal padding |

## Portrait Area

| Constant | Value | Description |
|----------|--------|-------------|
| `DialogPortraitWidth` | 120 | Portrait width in pixels |
| `DialogPortraitHeight` | 120 | Portrait height in pixels |
| `DialogPortraitX` | 30 | X position from widget left |
| `DialogPortraitY` | 30 | Y position from widget top |
| `DialogPortraitBorder` | 2 | Portrait border thickness |

## Text Area

| Constant | Value | Description |
|----------|--------|-------------|
| `DialogTextX` | 170 | Text area X position (after portrait) |
| `DialogTextY` | 30 | Text area Y position |
| `DialogTextWidth` | 500 | Text area width |
| `DialogTextHeight` | 180 | Text area height |
| `DialogTextLineHeight` | 18 | Height per text line |
| `DialogTextPadding` | 10 | Text area internal padding |

## Speaker Name Area

| Constant | Value | Description |
|----------|--------|-------------|
| `DialogSpeakerNameX` | 170 | Speaker name X position |
| `DialogSpeakerNameY` | 35 | Speaker name Y position |
| `DialogSpeakerNameHeight` | 25 | Height reserved for speaker name |

## Dialog Text Content

| Constant | Value | Description |
|----------|--------|-------------|
| `DialogContentY` | 65 | Dialog content Y position (after speaker name) |
| `DialogContentHeight` | 140 | Height available for dialog content |
| `DialogMaxLines` | 7 | Maximum visible text lines |

## Choice Buttons Area

| Constant | Value | Description |
|----------|--------|-------------|
| `DialogChoicesY` | 230 | Choices area Y position |
| `DialogChoicesHeight` | 140 | Height available for choices |
| `DialogChoiceHeight` | 25 | Height per choice button |
| `DialogChoiceSpacing` | 5 | Space between choice buttons |
| `DialogChoicePadding` | 8 | Internal padding for choice buttons |
| `DialogMaxVisibleChoices` | 5 | Maximum visible choices before scrolling |

## Typewriter Effect

| Constant | Value | Description |
|----------|--------|-------------|
| `DialogTypewriterSpeed` | 50 | Milliseconds per character |
| `DialogTypewriterFastSpeed` | 20 | Fast mode milliseconds per character |
| `DialogAutoAdvanceDelay` | 3000 | Auto-advance delay in milliseconds |

## Color Constants (RGBA values)

### Background Colors
| Color Component | R | G | B | A | Description |
|----------------|---|---|---|---|-------------|
| `DialogWidgetBackgroundR` | 25 | 25 | 35 | 240 | Widget background |
| `DialogWidgetBackgroundG` | 25 | 25 | 35 | 240 | Widget background |
| `DialogWidgetBackgroundB` | 35 | 25 | 35 | 240 | Widget background |
| `DialogWidgetBackgroundA` | 240 | 25 | 35 | 240 | Widget background |

### Border Colors
| Color Component | R | G | B | A | Description |
|----------------|---|---|---|---|-------------|
| `DialogWidgetBorderR` | 100 | 100 | 120 | 255 | Widget border |
| `DialogWidgetBorderG` | 100 | 100 | 120 | 255 | Widget border |
| `DialogWidgetBorderB` | 120 | 100 | 120 | 255 | Widget border |
| `DialogWidgetBorderA` | 255 | 100 | 120 | 255 | Widget border |

### Text Colors
| Color Component | R | G | B | A | Description |
|----------------|---|---|---|---|-------------|
| `DialogSpeakerNameR` | 255 | 220 | 100 | 255 | Speaker name text |
| `DialogSpeakerNameG` | 220 | 220 | 100 | 255 | Speaker name text |
| `DialogSpeakerNameB` | 100 | 220 | 100 | 255 | Speaker name text |
| `DialogSpeakerNameA` | 255 | 220 | 100 | 255 | Speaker name text |
| `DialogContentTextR` | 240 | 240 | 240 | 255 | Dialog content text |
| `DialogContentTextG` | 240 | 240 | 240 | 255 | Dialog content text |
| `DialogContentTextB` | 240 | 240 | 240 | 255 | Dialog content text |
| `DialogContentTextA` | 255 | 240 | 240 | 255 | Dialog content text |

### Choice Button Colors
| Color Component | R | G | B | A | Description |
|----------------|---|---|---|---|-------------|
| `DialogChoiceNormalR` | 60 | 60 | 80 | 255 | Normal choice background |
| `DialogChoiceNormalG` | 60 | 60 | 80 | 255 | Normal choice background |
| `DialogChoiceNormalB` | 80 | 60 | 80 | 255 | Normal choice background |
| `DialogChoiceNormalA` | 255 | 60 | 80 | 255 | Normal choice background |
| `DialogChoiceSelectedR` | 100 | 150 | 200 | 255 | Selected choice background |
| `DialogChoiceSelectedG` | 150 | 150 | 200 | 255 | Selected choice background |
| `DialogChoiceSelectedB` | 200 | 150 | 200 | 255 | Selected choice background |
| `DialogChoiceSelectedA` | 255 | 150 | 200 | 255 | Selected choice background |
| `DialogChoiceTextR` | 255 | 255 | 255 | 255 | Choice button text |
| `DialogChoiceTextG` | 255 | 255 | 255 | 255 | Choice button text |
| `DialogChoiceTextB` | 255 | 255 | 255 | 255 | Choice button text |
| `DialogChoiceTextA` | 255 | 255 | 255 | 255 | Choice button text |

### Portrait Colors
| Color Component | R | G | B | A | Description |
|----------------|---|---|---|---|-------------|
| `DialogPortraitBorderR` | 150 | 150 | 150 | 255 | Portrait border |
| `DialogPortraitBorderG` | 150 | 150 | 150 | 255 | Portrait border |
| `DialogPortraitBorderB` | 150 | 150 | 150 | 255 | Portrait border |
| `DialogPortraitBorderA` | 255 | 150 | 150 | 255 | Portrait border |
| `DialogPortraitPlaceholderR` | 80 | 80 | 100 | 255 | Portrait placeholder |
| `DialogPortraitPlaceholderG` | 80 | 80 | 100 | 255 | Portrait placeholder |
| `DialogPortraitPlaceholderB` | 100 | 80 | 100 | 255 | Portrait placeholder |
| `DialogPortraitPlaceholderA` | 255 | 80 | 100 | 255 | Portrait placeholder |

### Shadow Colors
| Color Component | R | G | B | A | Description |
|----------------|---|---|---|---|-------------|
| `DialogWidgetShadowR` | 0 | 0 | 0 | 120 | Drop shadow |
| `DialogWidgetShadowG` | 0 | 0 | 0 | 120 | Drop shadow |
| `DialogWidgetShadowB` | 0 | 0 | 0 | 120 | Drop shadow |
| `DialogWidgetShadowA` | 120 | 0 | 0 | 120 | Drop shadow |

## Widget States

| State | Description |
|-------|-------------|
| `DialogStateHidden` | Dialog widget is not visible |
| `DialogStateTyping` | Typewriter effect is displaying text |
| `DialogStateWaiting` | Waiting for user input to continue |
| `DialogStateChoices` | Displaying choices for user selection |
| `DialogStateAutoAdvance` | Auto-advancing to next dialog node |

## Input Controls

| Control | Action | Description |
|---------|--------|-------------|
| **Space** | Continue/Skip | Continue to next dialog or skip typewriter |
| **Enter** | Select | Select highlighted choice or continue |
| **Escape** | Cancel/Close | Cancel dialog or close widget |
| **Up/Down Arrows** | Navigate | Navigate through available choices |
| **Tab** | Fast Mode | Toggle typewriter fast mode |

## Implementation Notes

### Dialog Script Integration
- Dialog content loaded from external JSON files in `assets/dialogs/`
- Support for conditional branching based on game variables
- Character portrait integration with game assets
- Dynamic text replacement with game state variables

### Visual Design
- Portrait on left side with character image
- Speaker name prominently displayed
- Main dialog text with typewriter effect
- Choice buttons arranged vertically with clear selection highlighting
- Consistent with existing widget styling patterns

### Performance Considerations
- Text rendering optimized for frequent updates during typewriter effect
- Portrait caching to avoid repeated image loading
- Efficient choice button rendering for long dialog trees

This constants file provides the foundation for implementing the Dialog Widget following the established patterns from other UI components.