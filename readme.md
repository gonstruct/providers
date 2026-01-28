# Providers

A collection of reusable providers and adapters for external services such as encryption, mail, and storage. The goal is to provide simple, consistent interfaces that allow applications to interact with external systems without tight coupling to specific implementations.

Keep in mind that this is a work in progress, although ready for production use, it is not yet feature complete. It takes a lot to build a good provider/adapter, and is a continuous effort to make the adapters better and easier to use.

## Reason for this package

It was a huge pain to switch from one mail provider to another or to even compose a good mail message. This idea and base of the package is heavily inspired by the [Laravel](https://laravel.com) framework. However as we keep developing we probably will diverge more and more from it, as we add more features and make it more idiomatic to Go.

I know globals are looked down upon in Go, but I like them. This package is designed to make developer experience better when writing general API's in Go. It's fast, and easy to use. The testing and mocking capabilities are built-in so you can focus on writing your application logic instead of worrying about the implementation details of external services.

## ✨ Features

- Common interfaces for external services (mail, storage, encryption).
- Pluggable implementations (swap SMTP for a mock, or local storage for S3).
- Easy to test — use provided fakes to assert interactions in unit tests.
- Written in idiomatic Go, with simplicity and clarity in mind. (not entirely true, it uses globals which is looked down upon in Go, but I like it. I maybe provide better interfaces to use in something like a DI container like [https://github.com/samber/do](https://github.com/samber/do))

## 📦 Installation

```bash
go get github.com/gonstruct/providers
```

## 🚀 Usage

Examples can be found in the [examples](.examples) directory.

## 🧪 Testing

Each provider includes a `Fake()` function for unit tests that sets up a mock adapter and returns it for assertions:

```go
func TestEncryptData(t *testing.T) {
    fake := encryption.Fake()

    encrypted, _ := encryption.Encrypt("secret data")
    
    fake.AssertEncrypted(t, "secret data")
    fake.AssertEncryptCount(t, 1)
}

func TestSendEmail(t *testing.T) {
    fake := mail.Fake()

    mail.Send(myMailable)
    
    fake.AssertSentTo(t, "user@example.com")
    fake.AssertSentCount(t, 1)
}

func TestStoreFile(t *testing.T) {
    fake := storage.Fake()

    storage.PutFile("uploads", file.FromBytes("doc.pdf", content))
    
    fake.AssertStored(t, "uploads/doc.pdf")
}
```

> **Note:** Since providers use global state, `t.Parallel()` will not work correctly.

## 🛠️ Roadmap / Todo

- [ ] Improve documentation and examples 
- [ ] Change the mail provider
- [ ] Add caching provider with adapters for Redis, Memcached, etc.

## 🤝 Contributing

Contributions are welcome! If you'd like to add a new adapter or improve an existing one, feel free to open an issue or submit a PR.

## 📄 License

MIT License. See [LICENSE](LICENSE) for details.
