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

    Prefix_Filter<TC_shortcut> foo = FilterAPI<Prefix_Filter<TC_shortcut>>::ConstructFromAddCount(capacity);

    RedisModule_ModuleTypeSetValue(key, PFType, &foo);

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
