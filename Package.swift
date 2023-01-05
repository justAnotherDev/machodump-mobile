// swift-tools-version:5.3
import PackageDescription

let package = Package(
    name: "Macho",
    platforms: [
        .iOS(.v13)
    ],
    products: [
        .library(
            name: "Macho",
            targets: ["Macho"])
    ],
    dependencies: [],
    targets: [
        .binaryTarget(
            name: "Macho",
            url: "https://github.com/justAnotherDev/machodump-mobile/releases/download/1.0.0/Macho.xcframework.zip",
            checksum: "7301b7bef9e9fb6efb37ea6f379d8775cbff8dcad6e3f8dc2f99dc0adb9ecb6b"
        )
    ]
)