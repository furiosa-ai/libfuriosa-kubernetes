### XML topology 
The XML topology `test.xml` contains following topology for testing.

```text
Machine
├── Package (CPU)
│   ├── Host Bridge (Root Complex)
│   │   └── PCI Bridge
│   │       └── PCI Bridge
│   │           └── PCI Bridge
│   │               └── PCI Bridge
│   │                   └── PCI Bridge
│   │                       └── PCI Bridge
│   │                           └── PCI Bridge
│   │                               ├── NPU0(0000:27:00.0)
│   │                               └── NPU1(0000:2a:00.0)
│   └── Host Bridge (Root Complex)
│       └── PCI Bridge
│           └── PCI Bridge
│               └── PCI Bridge
│                   └── PCI Bridge
│                       └── PCI Bridge
│                           └── PCI Bridge
│                               └── PCI Bridge
│                                   ├── NPU2(0000:51:00.0)
│                                   └── NPU3(0000:57:00.0)
└── Package (CPU)
    ├── Host Bridge (Root Complex)
    │   └── PCI Bridge
    │       └── PCI Bridge
    │           └── PCI Bridge
    │               └── PCI Bridge
    │                   └── PCI Bridge
    │                       └── PCI Bridge
    │                           └── PCI Bridge
    │                               ├── NPU4(0000:9e:00.0)
    │                               └── NPU5(0000:a4:00.0)
    └── Host Bridge (Root Complex)
        └── PCI Bridge
            └── PCI Bridge
                └── PCI Bridge
                    └── PCI Bridge
                        └── PCI Bridge
                            └── PCI Bridge
                                └── PCI Bridge
                                    ├── NPU6(0000:c7:00.0)
                                    └── NPU7(0000:ca:00.0)
```

