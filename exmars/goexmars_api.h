#ifndef GOEXMARS_API_H
#define GOEXMARS_API_H

#ifdef __cplusplus
extern "C" {
#endif

typedef struct goexmars_fight_cfg_st {
	int coresize;
	int cycles;
	int maxprocess;
	int rounds;
	int maxwarriorlen;
	int minsep;
	int pspacesize;
	int fixpos;
} goexmars_fight_cfg_t;

void fight_1(char* w1, goexmars_fight_cfg_t* cfg, int* wins, int winsLen, int* ties, char* diagBuf, int diagCap, int* diagLen);
void fight_2(char* w1, char* w2, goexmars_fight_cfg_t* cfg, int* wins, int winsLen, int* ties, char* diagBuf, int diagCap, int* diagLen);
void fight_3(char* w1, char* w2, char* w3, goexmars_fight_cfg_t* cfg, int* wins, int winsLen, int* ties, char* diagBuf, int diagCap, int* diagLen);
void fight_4(char* w1, char* w2, char* w3, char* w4, goexmars_fight_cfg_t* cfg, int* wins, int winsLen, int* ties, char* diagBuf, int diagCap, int* diagLen);
void fight_5(char* w1, char* w2, char* w3, char* w4, char* w5, goexmars_fight_cfg_t* cfg, int* wins, int winsLen, int* ties, char* diagBuf, int diagCap, int* diagLen);
void fight_6(char* w1, char* w2, char* w3, char* w4, char* w5, char* w6, goexmars_fight_cfg_t* cfg, int* wins, int winsLen, int* ties, char* diagBuf, int diagCap, int* diagLen);

#ifdef __cplusplus
}
#endif

#endif
