{
  "coin": {
    "name": "Sinovate",
    "shortcut": "SIN",
    "label": "Sinovate",
    "alias": "sin"
  },
  "ports": {
    "backend_rpc": 20971,
    "backend_message_queue": 29000,
    "blockbook_internal": 9099,
    "blockbook_public": 9199
  },
  "ipc": {
    "rpc_url_template": "http://127.0.0.1:{{.Ports.BackendRPC}}",
    "rpc_user": "bitcoinrpc",
    "rpc_pass": "7RTk2adru5feH29JFiBVuCrJr",
    "rpc_timeout": 25,
    "message_queue_binding_template": "tcp://127.0.0.1:{{.Ports.BackendMessageQueue}}"
  },
  "backend": {
    "package_name": "backend-sinovate",
    "package_revision": "sinteam-1",
    "system_user": "sinteam",
    "version": "24.1.1",
    "binary_url": "https://github.com/SINOVATEblockchain/sinovate/releases/download/v24.1.1/sinovate-v24.1.1-x64-linux-gnu-with-cli-tools.tar.gz",
    "verification_type": "sha256",
    "verification_source": "28d8057bb92f390a792137aa7b3ab3aa6f77e2bbd17482a5237fbdd3a2775f1d",
    "extract_command": "bsdtar -C backend -xf",
    "exclude_files": [
      "sin-qt"
    ],
    "exec_command_template": "{{.Env.BackendInstallPath}}/{{.Coin.Alias}}/sind -datadir={{.Env.BackendDataPath}}/{{.Coin.Alias}}/backend -conf={{.Env.BackendInstallPath}}/{{.Coin.Alias}}/{{.Coin.Alias}}.conf -pid=/run/{{.Coin.Alias}}/{{.Coin.Alias}}.pid",
    "logrotate_files_template": "{{.Env.BackendDataPath}}/{{.Coin.Alias}}/backend/*.log",
    "postinst_script_template": "",
    "service_type": "forking",
    "service_additional_params_template": "",
    "protect_memory": true,
    "mainnet": true,
    "server_config_file": "bitcoin_like.conf",
    "client_config_file": "bitcoin_like_client.conf",
    "additional_params": {
      "deprecatedrpc": "estimatefee"
    }
  }
}
