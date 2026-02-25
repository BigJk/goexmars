#ifndef GOEXMARS_API_H
#define GOEXMARS_API_H

#ifdef __cplusplus
extern "C" {
#endif

void Fight2Warriors(char* w1, char* w2, int coresize, int cycles, int maxprocess, int rounds, int maxwarriorlen, int minsep, int pspacesize, int fixpos, int* win1, int* win2, int* equal);
void Fight1Warrior(char* w1, int coresize, int cycles, int maxprocess, int rounds, int maxwarriorlen, int minsep, int pspacesize, int fixpos, int* win1, int* win2, int* equal);

#ifdef __cplusplus
}
#endif

#endif
