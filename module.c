#include <bitset>
#include <functional>
#include <string>
#include <vector>

#include <stdlib.h>
#include "redismodule.h"
#include "Prefix-Filter/Tests/wrappers.hpp"

#ifdef __cplusplus
extern "C" {
#endif

static RedisModuleType * PFType;

static int pfGetValue(RedisModuleKey *key, RedisModuleType *expType, void **sbout) {
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


static int pfGetObject(RedisModuleKey *key, Prefix_Filter<SimdBlockFilterFixed<>> **sbout) {
    int i = pfGetValue(key, PFType, (void **)sbout);
    return i;
}

  static int PFExists_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc) {
    bool err = false;
    RedisModule_AutoMemory(ctx);
    std::hash<std::string> stdhash;

    if (argc < 2) {
      return RedisModule_WrongArity(ctx);
    }

    RedisModuleKey *key = RedisModule_OpenKey(ctx, argv[1], REDISMODULE_READ | REDISMODULE_WRITE);

    Prefix_Filter<SimdBlockFilterFixed<>> *pf = NULL;
    if (pfGetObject(key, &pf) < 0) {
      return RedisModule_ReplyWithError(ctx, "error fetching table by key");
    }

    if (pf == NULL) {
      return RedisModule_ReplyWithError(ctx, "internal server error");
    }

    size_t len = 0;

    const char * command = RedisModule_StringPtrLen(argv[0], &len);

    if (argc == 3 && ( strcasecmp(command, "PF.EXISTS") == 0)) {
      std::string str = RedisModule_StringPtrLen(argv[2], &len);

      unsigned long long h = stdhash(str);

      int result = FilterAPI<Prefix_Filter<SimdBlockFilterFixed<>>>::Contain(h, pf);

      RedisModule_ReplyWithLongLong(ctx, result);
      return REDISMODULE_OK;
    } else {
      int * successes = NULL;
      u64 * foos = NULL;
      unsigned long i = 0;
      unsigned long foolen = argc - 2;
      foos = (unsigned long *)malloc(sizeof(*foos)*(foolen));
      if (foos == NULL) {
        exit(-5);
      }
      for (i = 0; i < foolen; i++) {
        foos[i] = stdhash(RedisModule_StringPtrLen(argv[i+2], &len));
      }
      successes = FilterAPI<Prefix_Filter<SimdBlockFilterFixed<>>>::MultiExists(foolen, foos, pf);
      if (successes == NULL) {
        err = true;
        goto cleanup;
      }

    cleanup:
      if (foos) free(foos);
      if (err) {
        if (successes) free(successes);
        return REDISMODULE_ERR;
      } else {
        RedisModule_ReplyWithArray(ctx, foolen);
        for (i = 0; i < foolen; i++) {
          RedisModule_ReplyWithLongLong(ctx, successes[i]);
        }
        free(successes);
        return REDISMODULE_OK;
      }

    }

    return REDISMODULE_OK;
  }

static int PFInfo_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc) {
    RedisModule_AutoMemory(ctx);
    if (argc != 2) {
        return RedisModule_WrongArity(ctx);
    }

    RedisModuleKey *key = RedisModule_OpenKey(ctx, argv[1], REDISMODULE_READ);
    Prefix_Filter<SimdBlockFilterFixed<>> *pf = NULL;
    if (pfGetObject(key, &pf) < 0) {
        return RedisModule_ReplyWithError(ctx, "error fetching table by key");
    }

    RedisModule_ReplyWithArray(ctx, 2 * 2);
    RedisModule_ReplyWithSimpleString(ctx, "Capacity");
    RedisModule_ReplyWithLongLong(ctx, pf->get_max_capacity());
    RedisModule_ReplyWithSimpleString(ctx, "Size");
    RedisModule_ReplyWithLongLong(ctx, pf->get_cap());
    return REDISMODULE_OK;
}

static int PFAdd_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc) {
    bool err = false;
    RedisModule_AutoMemory(ctx);
    std::hash<std::string> stdhash;

    if (argc < 2) {
      return RedisModule_WrongArity(ctx);
    }

    RedisModuleKey *key = RedisModule_OpenKey(ctx, argv[1], REDISMODULE_READ | REDISMODULE_WRITE);

    Prefix_Filter<SimdBlockFilterFixed<>> *pf = NULL;
    if (pfGetObject(key, &pf) < 0) {
      return RedisModule_ReplyWithError(ctx, "error fetching table by key");
    }

    if (pf == NULL) {
      return RedisModule_ReplyWithError(ctx, "internal server error");
    }

    size_t len = 0;
    if (argc == 3) {
      std::string str = RedisModule_StringPtrLen(argv[2], &len);

      unsigned long long h = stdhash(str);

      FilterAPI<Prefix_Filter<SimdBlockFilterFixed<>>>::Add(h, pf);

      RedisModule_ReplyWithSimpleString(ctx, "OK");

      return REDISMODULE_OK;
    } else {
      int * successes = NULL;
      u64 * foos = NULL;
      unsigned long i = 0;
      unsigned long foolen = argc - 2;
      foos = (unsigned long *)malloc(sizeof(*foos)*(foolen));
      if (foos == NULL) {
        exit(-5);
      }
      for (i = 0; i < foolen; i++) {
        foos[i] = stdhash(RedisModule_StringPtrLen(argv[i+2], &len));
      }
      successes = FilterAPI<Prefix_Filter<SimdBlockFilterFixed<>>>::MultiAdd(foolen, foos, pf);
      if (successes == NULL) {
        err = true;
        goto cleanup;
      }

    cleanup:
      if (foos) free(foos);
      if (err) {
        if (successes) free(successes);
        return REDISMODULE_ERR;
      } else {
        RedisModule_ReplyWithArray(ctx, foolen);
        for (i = 0; i < foolen; i++) {
          RedisModule_ReplyWithLongLong(ctx, successes[i]);
        }
        free(successes);
        return REDISMODULE_OK;
      }
    }
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

    Prefix_Filter<SimdBlockFilterFixed<>>* table = new Prefix_Filter<SimdBlockFilterFixed<>>(capacity, loads);

    RedisModule_ModuleTypeSetValue(key, PFType, table);

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
    if (RedisModule_Init(ctx,"pf",1,REDISMODULE_APIVER_1) == REDISMODULE_ERR)
      return REDISMODULE_ERR;

    if (RedisModule_CreateCommand(ctx,"pf.reserve",
                                  PFReserve_RedisCommand, "write",
                                  1, 1, 1) == REDISMODULE_ERR)
      return REDISMODULE_ERR;

    if (RedisModule_CreateCommand(ctx,"pf.info",
                                  PFInfo_RedisCommand, "readonly",
                                  1, 1, 1) == REDISMODULE_ERR)
      return REDISMODULE_ERR;

    if (RedisModule_CreateCommand(ctx,"pf.add",
                                  PFAdd_RedisCommand, "write",
                                  1, 1, 1) == REDISMODULE_ERR)
      return REDISMODULE_ERR;

    if (RedisModule_CreateCommand(ctx,"pf.madd",
                                  PFAdd_RedisCommand, "write",
                                  1, 1, 1) == REDISMODULE_ERR)
      return REDISMODULE_ERR;

    if (RedisModule_CreateCommand(ctx,"pf.mexists",
                                  PFExists_RedisCommand, "readonly",
                                  1, 1, 1) == REDISMODULE_ERR)
      return REDISMODULE_ERR;

    if (RedisModule_CreateCommand(ctx,"pf.exists",
                                  PFExists_RedisCommand, "readonly",
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
