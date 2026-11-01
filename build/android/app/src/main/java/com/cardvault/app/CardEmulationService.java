package com.cardvault.app;

import android.nfc.cardemulation.HostApduService;
import android.os.Bundle;
import android.util.Log;

public class CardEmulationService extends HostApduService {
    private static final String TAG = "cardVaultHCE";
    private static final byte[] SUCCESS = new byte[] { (byte) 0x90, 0x00 };
    private static final byte[] NOT_SUPPORTED = new byte[] { (byte) 0x6A, (byte) 0x82 };

    @Override
    public byte[] processCommandApdu(byte[] commandApdu, Bundle extras) {
        if (commandApdu == null || commandApdu.length < 4) {
            return NOT_SUPPORTED;
        }
        if (commandApdu[1] == (byte) 0xA4) {
            Log.d(TAG, "SELECT received");
            return SUCCESS;
        }
        return SUCCESS;
    }

    @Override
    public void onDeactivated(int reason) {
        Log.d(TAG, "HCE deactivated: " + reason);
    }
}
