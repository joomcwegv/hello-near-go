# Hello NEAR — Go тілінде алғашқы смарт-келісімшарт 🚀

Бұл — **NEAR Protocol** блокчейнінде **Go (Golang)** тілінде жазылған алғашқы смарт-келісімшарт.

## 📦 Технологиялар
- **Go (Golang)** — смарт-келісімшарт тілі
- **TinyGo** — Go кодын WebAssembly (WASM) форматына аударады
- **NEAR Protocol** — блокчейн платформасы
- **near-cli-rs** — NEAR-мен терминал арқылы жұмыс жасайды

## 🛠 Орнату

```bash
# TinyGo орнату (Windows)
winget install TinyGo.TinyGo

# NEAR CLI орнату
npm install -g near-cli-rs@latest
```

## 🔨 Компиляция

```bash
tinygo build -size short -no-debug -panic=trap -scheduler=none -gc=leaking -o contract.wasm -target wasm-unknown ./
```

## 🚀 Testnet-ке деплой

```bash
# Аккаунт жасау
near account create-account sponsor-by-faucet-service YOUR_NAME.testnet autogenerate-new-keypair save-to-legacy-keychain network-config testnet create

# Деплой
near contract deploy YOUR_NAME.testnet use-file contract.wasm without-init-call network-config testnet sign-with-legacy-keychain send
```

## 📡 Функцияларды шақыру

```bash
# Хабарлама алу (тегін, газсыз)
near contract call-function as-read-only YOUR_NAME.testnet get_message json-args '{}' network-config testnet now
```

## 📋 Нәтиже

```
Logs:
│      Hello from Go on NEAR!

Function execution return value:
Hello from Go on NEAR!
```

## 🌐 Testnet Explorer
https://explorer.testnet.near.org/accounts/neargotest2025.testnet
