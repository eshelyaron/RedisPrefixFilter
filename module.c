#include <bitset>
#include <functional>
#include <string>
#include <vector>


#include "redismodule.h"
#include <stdlib.h>
#include "Prefix-Filter/Tests/wrappers.hpp"

#ifdef __cplusplus
extern "C" {
#endif

static RedisModuleType * PFType;

int HelloworldRand_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc) {
    RedisModule_ReplyWithLongLong(ctx,rand());
    FilterAPI<Prefix_Filter<TC_shortcut>>::ConstructFromAddCount(1000000);
    return REDISMODULE_OK;
}

static int getValue(RedisModuleKey *key, RedisModuleType *expType, void **sbout) {
    *sbout = NULL;
    if (key == NULL) {
        return -1;
    }
    int type = RedisModule_KeyType(key);
    if (type == REDISMODULE_KEYTYPE_EMPTY) {
        return -2;
    } else if (type == REDISMODULE_KEYTYPE_MODULE &&
               RedisModule_ModuleTypeGetType(key) == expType) {
        *sbout = RedisModule_ModuleTypeGetValue(key);
        return 0;
    } else {
        return -3;
    }
}


static int pfGetObject(RedisModuleKey *key, Prefix_Filter<TC_shortcut> **sbout) {
    printf("pfGetObject\n");
    int i = getValue(key, PFType, (void **)sbout);
    printf("%d\n", i);
    return i;
}

static int PFAdd_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc) {
    RedisModule_AutoMemory(ctx);
    std::hash<std::string> stdhash;

    if (argc < 2) {
        return RedisModule_WrongArity(ctx);
    }

    printf("got here\n");
    RedisModuleKey *key = RedisModule_OpenKey(ctx, argv[1], REDISMODULE_READ | REDISMODULE_WRITE);

    Prefix_Filter<TC_shortcut> *pf = NULL;
    if (pfGetObject(key, &pf) < 0) {
      return RedisModule_ReplyWithError(ctx, "error fetching table by key");
    }

    printf("got there\n");
    printf("%p\n", pf);
    if (pf == NULL) {
      return RedisModule_ReplyWithError(ctx, "internal server error");
    }

    size_t len = 0;
    std::string str = RedisModule_StringPtrLen(argv[2], &len);

    unsigned long long h = stdhash(str);

    printf("got here again\n");

    printf("size: %lu\n", FilterAPI<Prefix_Filter<TC_shortcut>>::get_cap(pf));

    FilterAPI<Prefix_Filter<TC_shortcut>>::Add(h, pf);

    RedisModule_ReplyWithSimpleString(ctx, "OK");

    printf("done\n");
    return REDISMODULE_OK;
}

static int PFReserve_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc) {
    RedisModule_AutoMemory(ctx);

    if (argc < 2) {
        return RedisModule_WrongArity(ctx);
    }

    long long capacity;
    if (RedisModule_StringToLongLong(argv[2], &capacity) != REDISMODULE_OK) {
        return RedisModule_ReplyWithError(ctx, "ERR bad capacity");
    } else if (capacity <= 0) {
        return RedisModule_ReplyWithError(ctx, "ERR (capacity should be larger than 0)");
    }

    RedisModuleKey *key = RedisModule_OpenKey(ctx, argv[1], REDISMODULE_READ | REDISMODULE_WRITE);
    float loads[2] = {.95, .95};

    Prefix_Filter<TC_shortcut>* foo = new Prefix_Filter<TC_shortcut>(capacity, loads);

    printf("%p\n", foo);
    printf("func: %u\n", FilterAPI<Prefix_Filter<TC_shortcut>>::get_functionality(foo));
    printf("cap: %lu\n", FilterAPI<Prefix_Filter<TC_shortcut>>::get_cap(foo));

    FilterAPI<Prefix_Filter<TC_shortcut>>::Add(9999, foo);

    printf("new cap: %lu\n", FilterAPI<Prefix_Filter<TC_shortcut>>::get_cap(foo));
    RedisModule_ModuleTypeSetValue(key, PFType, foo);

    RedisModule_ReplyWithSimpleString(ctx, "OK");

    return REDISMODULE_OK;
}


static void BFRdbSave(RedisModuleIO *io, void *obj) {}

static void *BFRdbLoad(RedisModuleIO *io, int encver) {
  return NULL;
}

static void BFAofRewrite(RedisModuleIO *aof, RedisModuleString *key, void *value) {
  (void)value;
}

static void BFFree(void *value) { (void)value; }

int RedisModule_OnLoad(RedisModuleCtx *ctx, RedisModuleString **argv, int argc) {
    if (RedisModule_Init(ctx,"avromodule",1,REDISMODULE_APIVER_1)
        == REDISMODULE_ERR) return REDISMODULE_ERR;

    if (RedisModule_CreateCommand(ctx,"pf.reserve",
                                  PFReserve_RedisCommand, "write",
                                  1, 1, 1) == REDISMODULE_ERR)
      return REDISMODULE_ERR;

    if (RedisModule_CreateCommand(ctx,"pf.add",
                                  PFAdd_RedisCommand, "write",
                                  1, 1, 1) == REDISMODULE_ERR)
      return REDISMODULE_ERR;

    static RedisModuleTypeMethods typeprocs = {.version = REDISMODULE_TYPE_METHOD_VERSION,
                                               .rdb_load = BFRdbLoad,
                                               .rdb_save = BFRdbSave,
                                               .aof_rewrite = BFAofRewrite,
                                               .free = BFFree};

    PFType = RedisModule_CreateDataType(ctx, "PFilter--", 0, &typeprocs);

    return REDISMODULE_OK;
}

#ifdef __cplusplus
}
#endif
