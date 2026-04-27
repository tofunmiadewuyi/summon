#include <stdint.h>

void clearTapBindings(void);
void addTapBinding(uint16_t keycode, uint32_t modifiers, const char* appName);
int  startEventTap(void);
int  isTapRunning(void);
int  captureNextCombo(uint16_t* keycode, uint64_t* modifiers);
