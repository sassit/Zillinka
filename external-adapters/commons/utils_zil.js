const { Zilliqa } = require('@zilliqa-js/zilliqa');
const { bytes } = require('@zilliqa-js/util');

function setup_chain_and_wallet(/*bool*/testnet)
{
  let zilliqa_chain = null;
  let chainId = 0;
  let privateKeys = '';
  if (testnet) {
    zilliqa_chain = new Zilliqa('https://dev-api.zilliqa.com');
    chainId = 333;
    privateKey =  // corresponding address: 0x56A7812f68cbF83194a4a777D2310Aa7A378C9D8
      'b74501e0d2d047e8aaa2353020b46f31d396f92d05843665573300995e3aed88';
    uxt_oracle_addr = '0xaa28674d160a7b74cd6b3eedce733ac5c01cd26a';
    uxt_oracle_client_addr = '0x61b0e75ad884ccb3253c9eaa6396dea3c80f2f67';
    rhine_oracle_addr = '0xbd0a71b5490291cf99b00a733068627abdac6fcc';
    rhine_oracle_client_addr = '0xc9357c221cebf5f20e35153103032747a0d99bbf';
  }
  else { // Isolated server / Simulated ENV
    zilliqa_chain = new Zilliqa('https://zilliqa-isolated-server.zilliqa.com/');
    chainId = 222;
    privateKey = // corresponding address: 0x0427fd49d99d8eBDdC7827ECa98AC01D19601E5d
    'c0f1b81c2b541fed2d75d9eae6e096fa1b74b0ec0add48e1543e70ebfaaeed99';
    uxt_oracle_addr = '0x827c1f98a934de858f875c0d7a489a24a1d119ed';
    uxt_oracle_client_addr = '0xb80ad4de4ace27c7313d2439d8615ce7bad9a23d';
    rhine_oracle_addr = '0x7eabaf2ac5b9a415e27c1e78cf9d831e46446f5b';
    rhine_oracle_client_addr = '0x40ec20eaea9a1345d548b5490da50a9437b9c800';
  }
  const msgVersion = 1; // current msgVersion
  const VERSION = bytes.pack(chainId, msgVersion);
  return {"zilliqa": zilliqa_chain,
          "VERSION": VERSION,
          "privateKey": privateKey,
          "addresses": {
                "UnixTimeOracle": uxt_oracle_addr,
                "UnixTimeOracleClient": uxt_oracle_client_addr,
                "RhineGaugeOracle": rhine_oracle_addr,
                "RhineGaugeOracleClient": rhine_oracle_client_addr,
              },
            };
}


// deploy a contract (in a string) given init parameters, blockchain and tx txParams
async function deploy_contract(/*string*/sc_string, /*JSON*/init, /*JSON*/ bc_setup, /*JSON*/ tx_settings)
{
  const contract = bc_setup.zilliqa.contracts.new(sc_string, init);
  const [tx, sc] = await contract.deploy(
    { version: bc_setup.VERSION, gasPrice: tx_settings.gas_price, gasLimit: tx_settings.gas_limit, },
    tx_settings.attempts, 1000, true,
  );
  return [tx, sc];
}

// call a smart contract's transition with given args and an amount to send from a given public key
async function call_contract( /*contract*/sc, /*string*/transition_name, /*array*/args,
                        /*BN*/amt_as_BN, /*string*/caller_pub_key, /*JSON*/ bc_setup, /*JSON*/tx_settings)
{
   const tx = await sc.call(
     transition_name,
     args,
     { version: bc_setup.VERSION, amount: amt_as_BN, gasPrice: tx_settings.gas_price,
       gasLimit: tx_settings.gas_limit, pubKey: caller_pub_key, },
     tx_settings.attempts, 1000, true,
   );
   return tx;
}


exports.setup_chain_and_wallet = setup_chain_and_wallet;
exports.deploy_contract = deploy_contract;
exports.call_contract = call_contract;
