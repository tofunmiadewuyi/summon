#include <stdint.h>

void clearTapBindings(void);
void addTapBinding(uint16_t keycode, uint32_t modifiers, const char* appName);
void startEventTap(void);
int  captureNextCombo(uint16_t* keycode, uint64_t* modifiers);
