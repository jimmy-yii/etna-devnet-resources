import { bytesToHex, hexToBytes } from '@noble/hashes/utils';
import fs from "fs"
import { getPrivateKeyWallet } from "../lib/wallet";
import { secp256k1 } from 'avalanchejs-bleeding-edge';


let privateKeyBytes: Uint8Array
if (fs.existsSync("privateKey.txt")) {
    privateKeyBytes = hexToBytes(fs.readFileSync("privateKey.txt", "utf8").trim());
} else {
    privateKeyBytes = secp256k1.randomPrivateKey();
    fs.writeFileSync("privateKey.txt", bytesToHex(privateKeyBytes));
}

const wallet = getPrivateKeyWallet(privateKeyBytes);

console.log(await wallet.getAddress());
