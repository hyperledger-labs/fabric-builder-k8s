/*
 * SPDX-License-Identifier: Apache-2.0
 */

import {Context, Contract, Transaction} from 'fabric-contract-api';
import {TextDecoder, TextEncoder} from 'util';

const encoder = new TextEncoder();
const decoder = new TextDecoder();

class SampleContract extends Contract {
  @Transaction()
  async PutValue(ctx: Context, key: string, value: string) {
    await ctx.stub.putState(key, encoder.encode(value));
  }

  @Transaction()
  async GetValue(ctx: Context, key: string) {
    const value = await ctx.stub.getState(key);
    return decoder.decode(value);
  }
}

exports.contracts = [SampleContract];
