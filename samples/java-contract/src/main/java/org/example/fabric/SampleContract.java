/*
 * SPDX-License-Identifier: Apache-2.0
 */

package org.example.fabric;

import static java.nio.charset.StandardCharsets.UTF_8;

import org.hyperledger.fabric.contract.Context;
import org.hyperledger.fabric.contract.ContractInterface;
import org.hyperledger.fabric.contract.annotation.Contract;
import org.hyperledger.fabric.contract.annotation.Default;
import org.hyperledger.fabric.contract.annotation.Transaction;

@Contract(name = "sample")
@Default
public final class SampleContract implements ContractInterface {

    /**
     * Adds a key value pair to the world state.
     *
     * @param ctx the transaction context
     * @param key the key
     * @param value the value
     */
    @Transaction(intent = Transaction.TYPE.SUBMIT)
    public void PutValue(final Context ctx, final String key, final String value) {
        ctx.getStub().putState(key, value.getBytes());
    }

    /**
     * Gets the value for a key from the world state.
     *
     * @param ctx the transaction context
     * @param key the key
     * @return the value
     */
    @Transaction(intent = Transaction.TYPE.EVALUATE)
    public String GetValue(final Context ctx, final String key) {
        return new String(ctx.getStub().getState(key), UTF_8);
    }
}

