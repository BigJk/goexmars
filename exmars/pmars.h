#ifndef PMARS_H
#define PMARS_H

#include <stdio.h>
#include <ctype.h>
#include <string.h>
#include <time.h>
#include <stdio.h>
#include "pmars.h"
#include "exhaust.h"
#include "sim.h"
#include "insn_help.h"
#include "insn.h"

/* ***********************************************************************
	 function prototypes
	 *********************************************************************** */
int parse_param(int argc, char *argv[]);
int eval_expr(mars_t* mars, char *expr, long *result);
int assemble(char *fName, int aWarrior);
void disasm(mars_t* mars, mem_struct * cells, ADDR_T n, ADDR_T offset);
void simulator1(void);
char *locview(ADDR_T loc, char *outp);
int cdb(char *msg);
int score(int warnum);
int deaths(int warnum);
void results(FILE * outp);
void Exit(int code);
void reset_regs(mars_t* mars);

char *cellview(mars_t* mars, mem_struct * cell, char *outp);

char* pstrdup(char *);
char* pstrcat(char *, char *);
const char* pstrchr(const char *, int);
uChar ch_in_set(char c, const char* s);
uChar skip_space(char *, uChar);
uChar str_in_set(char *, const char *s[]);
int get_token(char * , uChar *, char *);
void to_upper(char *);


int assemble_warrior2(mars_t* mars, char* redstr, warrior_struct* w);
void Fight2Warriors(char*, char*, int, int, int, int, int, int, int, int, int*, int*, int*);
void Fight1Warrior(char*, int, int, int, int, int, int, int, int, int*, int*, int*);

/* ****************** required local prototypes ********************* */

static void textout(char *);
static void errprn(mars_t* mars, errType code, line_st* aline, const char* arg);

/* ***************** conforming local prototypes ******************** */
static int trav2(mars_t* mars, char* buffer, char* dest, int wdecl);
static int normalize(mars_t* mars, long);
static int blkfor(mars_t* mars, char *, char *);
static int equtbl(mars_t* mars, char *);
static int equsub(mars_t* mars, char *, char *, int, ref_st *);
static ref_st *lookup(mars_t* mars, char *);
static grp_st *addsym(mars_t* mars, char *, grp_st *);
static src_st *addlinesrc(mars_t* mars, char *, uShrt);
static void newtbl(mars_t* mars);
static void addpredef(mars_t* mars, char *, U32_T);
static void addpredefs(mars_t* mars);
static void addline(mars_t* mars, char*, src_st*, uShrt);
static void disposeline(line_st *), disposegrp(grp_st *);
static void disposetbl(ref_st *, ref_st *);
static void cleanmem(mars_t* mars);
static void nocmnt(char *);
static void automaton(mars_t* mars, char *, stateCol, mem_struct *);
static void dfashell(mars_t* mars, char *, mem_struct *);
static void expand(mars_t* mars, uShrt);

int denormalize(mars_t* mars, int x);
void disposesrc(src_st* r);

int globalswitch_warrior(mars_t* mars, warrior_struct* w, char* str, uShrt idx, uShrt loc, uShrt lspnt);
void encode_warrior(mars_t* mars, warrior_struct* w, uShrt sspnt);
s32_t rng(s32_t seed);
void panic(char *msg);
int assemble_warrior(mars_t* mars, char* fName, warrior_struct* w);
void clear_results(mars_t* mars);
void save_pspaces(mars_t* mars);
void amalgamate_pspaces(mars_t* mars);
void accumulate_results(mars_t* mars);
void check_sanity(mars_t* mars);
void readargs(int argc, char** argv, mars_t* mars);
void usage(void);
void load_warriors(mars_t* mars);
void set_starting_order(unsigned int round, mars_t* mars);
void output_results(mars_t* mars);
s32_t compute_positions(s32_t seed, mars_t* mars);
void npos(s32_t *seed, mars_t* mars);
int posit(s32_t *seed, mars_t* mars);
mars_t* init(int argc, char** argv);
void pmars2exhaust(mars_t* mars, warrior_struct** warriors, int wCount);

char* eval(mars_t* mars, int prevPrec, long val1, int oper1, char *expr, long *result);
char	 *getreg(mars_t* mars, char *expr, int regId, long *val);
char	 *getval(mars_t* mars, char *expr, long *val);
char	 *getop(char *expr, char *op);
long		calc(mars_t* mars, long x, long y, int op);

#endif
