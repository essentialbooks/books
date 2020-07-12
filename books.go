package main

var (
	bookGo = &Book{
		Title:          "Go",
		TitleLong:      "Essential Go",
		Dir:            "go",
		CoverImageName: "Go.png",
		// https://www.notion.so/2cab1ed2b7a44584b56b0d3ca9b80185
		NotionStartPageID: "2cab1ed2b7a44584b56b0d3ca9b80185",
	}
	bookCsharp = &Book{
		Title:          "C#",
		TitleLong:      "Essential C#",
		Dir:            "csharp",
		CoverImageName: "CSharp.png",
		// https://www.notion.so/kjkpublic/Essential-C-896da5248e65414ab645dd45985879a1
		NotionStartPageID: "896da5248e65414ab645dd45985879a1",
	}
	bookPython = &Book{
		Title:          "Python",
		TitleLong:      "Essential Python",
		Dir:            "python",
		CoverImageName: "Python.png",
		// https://www.notion.so/kjkpublic/Essential-Python-12e6f78e68a5497290c96e1365ae6259
		NotionStartPageID: "12e6f78e68a5497290c96e1365ae6259",
	}
	bookKotlin = &Book{
		Title:          "Kotlin",
		TitleLong:      "Essential Kotlin",
		Dir:            "kotlin",
		CoverImageName: "Kotlin.png",
		// https://www.notion.so/kjkpublic/Essential-Kotlin-2bdd47318f3a4e8681dda289a8b3472b
		NotionStartPageID: "2bdd47318f3a4e8681dda289a8b3472b",
	}
	bookJavaScript = &Book{
		Title:          "JavaScript",
		TitleLong:      "Essential JavaScript",
		Dir:            "javascript",
		CoverImageName: "JavaScript.png",
		// https://www.notion.so/kjkpublic/Essential-Javascript-0b121710a160402fa9fd4646b87bed99
		NotionStartPageID: "0b121710a160402fa9fd4646b87bed99",
	}
	bookDart = &Book{
		Title:          "Dart",
		TitleLong:      "Essential Dart",
		Dir:            "dart",
		CoverImageName: "Dart.png",
		// 	https://www.notion.so/kjkpublic/Essential-Dart-0e2d248bf94b4aebaefbcf51ae435df0
		NotionStartPageID: "0e2d248bf94b4aebaefbcf51ae435df0",
	}
	bookJava = &Book{
		Title:          "Java",
		TitleLong:      "Essential Java",
		Dir:            "java",
		CoverImageName: "Java.png",
		// https://www.notion.so/kjkpublic/Essential-Java-d37cda98a07046f6b2cc375731ea3bdb
		NotionStartPageID: "d37cda98a07046f6b2cc375731ea3bdb",
	}
	bookAndroid = &Book{
		Title:          "Android",
		TitleLong:      "Essential Android",
		Dir:            "android",
		CoverImageName: "Android.png",
		// https://www.notion.so/kjkpublic/Essential-Android-f90b0a6b648343e28dc5ed6e8f5c0780
		NotionStartPageID: "f90b0a6b648343e28dc5ed6e8f5c0780",
	}
	bookSql = &Book{
		Title:          "SQL",
		TitleLong:      "Essential SQL",
		Dir:            "sql",
		CoverImageName: "SQL.png",
		// https://www.notion.so/kjkpublic/Essential-SQL-d1c8bb39bad4494e80abe28414c3d80e
		NotionStartPageID: "d1c8bb39bad4494e80abe28414c3d80e",
	}
	bookCpp = &Book{
		Title:          "C++",
		TitleLong:      "Essential C++",
		Dir:            "cpp",
		CoverImageName: "Cpp.png",
		// https://www.notion.so/kjkpublic/Essential-C-ad527dc6d4a7420b923494d0b9bfb560
		NotionStartPageID: "ad527dc6d4a7420b923494d0b9bfb560",
	}
	bookIOS = &Book{
		Title:          "iOS",
		TitleLong:      "Essential iOS",
		Dir:            "ios",
		CoverImageName: "iOS.png",
		// https://www.notion.so/kjkpublic/Essential-iOS-3626edc1bd044431afddc89648a7050f
		NotionStartPageID: "3626edc1bd044431afddc89648a7050f",
	}
	bookPostgresql = &Book{
		Title:          "PostgreSQL",
		TitleLong:      "Essential PostgreSQL",
		Dir:            "postgresql",
		CoverImageName: "PostgreSQL.png",
		// https://www.notion.so/kjkpublic/Essential-PostgreSQL-799304340f2c4081b6c4b7eb28df368e
		NotionStartPageID: "799304340f2c4081b6c4b7eb28df368e",
	}
	bookMysql = &Book{
		Title:          "MySQL",
		TitleLong:      "Essential MySQL",
		Dir:            "mysql",
		CoverImageName: "MySQL.png",
		// https://www.notion.so/kjkpublic/Essential-MySQL-4489ab73989f4ae9912486561e165deb
		NotionStartPageID: "4489ab73989f4ae9912486561e165deb",
	}
	bookAlgorithm = &Book{
		Title:          "Algorithms",
		TitleLong:      "Essential Algorithms",
		Dir:            "algorithms",
		CoverImageName: "Algorithms.png",
		// https://www.notion.so/kjkpublic/Essential-Algorithms-039ec42ee62f412e983e6d5b6b201b60
		NotionStartPageID: "039ec42ee62f412e983e6d5b6b201b60",
	}
	bookBash = &Book{
		Title:          "Bash",
		TitleLong:      "Essential Bash",
		Dir:            "bash",
		CoverImageName: "Bash.png",
		// https://www.notion.so/kjkpublic/Essential-Bash-77d28932012b489db9a6d0b349cea865
		NotionStartPageID: "77d28932012b489db9a6d0b349cea865",
	}
	bookC = &Book{
		Title:          "C",
		TitleLong:      "Essential C",
		Dir:            "c",
		CoverImageName: "C.png",
		// https://www.notion.so/kjkpublic/Essential-C-84ae4145718e4b7b8cb43cf10ee4db6a
		NotionStartPageID: "84ae4145718e4b7b8cb43cf10ee4db6a",
	}
	bookCSS = &Book{
		Title:          "CSS",
		TitleLong:      "Essential CSS",
		Dir:            "css",
		CoverImageName: "CSS.png",
		// https://www.notion.so/kjkpublic/Essential-CSS-18bfe038109649f48904e71ca18d76ed
		NotionStartPageID: "18bfe038109649f48904e71ca18d76ed",
	}
	bookGit = &Book{
		Title:          "Git",
		TitleLong:      "Essential Git",
		Dir:            "git",
		CoverImageName: "Git.png",
		// https://www.notion.so/kjkpublic/Essential-Git-37913107a4194981b1fc745928c0df66
		NotionStartPageID: "37913107a4194981b1fc745928c0df66",
	}
	bookHTML = &Book{
		Title:          "HTML",
		TitleLong:      "Essential HTML",
		Dir:            "html",
		CoverImageName: "HTML.png",
		// https://www.notion.so/kjkpublic/Essential-HTML-1c13e594ccd5472fb58d4c56379e7540
		NotionStartPageID: "1c13e594ccd5472fb58d4c56379e7540",
	}
	bookHTMLCanvas = &Book{
		Title:          "HTML Canvas",
		TitleLong:      "Essential HTML Canvas",
		Dir:            "htmlcanvas",
		CoverImageName: "HTMLCanvas.png",
		// https://www.notion.so/kjkpublic/Essential-HTML5-Canvas-227fa77d624c441d98011d7c998609a6
		NotionStartPageID: "227fa77d624c441d98011d7c998609a6",
	}
	bookNETFramework = &Book{
		Title:          "NET Framework",
		TitleLong:      "Essential NET Framework",
		Dir:            "netframework",
		CoverImageName: "NETFramework.png",
		// https://www.notion.so/kjkpublic/Essential-NET-framework-6289ba4610874393844dfbae890ed035
		NotionStartPageID: "6289ba4610874393844dfbae890ed035",
	}
	bookNode = &Book{
		Title:          "Node.js",
		TitleLong:      "Essential Node.js",
		Dir:            "nodejs",
		CoverImageName: "Node.js.png",
		// https://www.notion.so/kjkpublic/Essential-Node-a44f25f371164e69a70578bd98e71eb1
		NotionStartPageID: "a44f25f371164e69a70578bd98e71eb1",
	}
	bookObjectiveC = &Book{
		Title:          "Objective-C",
		TitleLong:      "Essential Objective-C",
		Dir:            "objectivec",
		CoverImageName: "ObjectiveC.png",
		// https://www.notion.so/kjkpublic/Essential-Objective-C-d16f98c8a3d641cd88a65fb67e6b0081
		NotionStartPageID: "d16f98c8a3d641cd88a65fb67e6b0081",
	}
	bookPHP = &Book{
		Title:          "PHP",
		TitleLong:      "Essential PHP",
		Dir:            "php",
		CoverImageName: "PHP.png",
		// https://www.notion.so/kjkpublic/Essential-PHP-b64d2e8d06e04dbea37e6e7d0e06bb48
		NotionStartPageID: "b64d2e8d06e04dbea37e6e7d0e06bb48",
	}
	bookPowershell = &Book{
		Title:          "PowerShell",
		TitleLong:      "Essential PowerShell",
		Dir:            "powershell",
		CoverImageName: "PowerShell.png",
		// https://www.notion.so/kjkpublic/Essential-Powershell-6042c6d3aed54250a12900e7f6b326e0
		NotionStartPageID: "6042c6d3aed54250a12900e7f6b326e0",
	}
	bookReact = &Book{
		Title:          "React",
		TitleLong:      "Essential React",
		Dir:            "react",
		CoverImageName: "Ruby.png", // TODO: need the cover
		// https://www.notion.so/kjkpublic/Essential-React-2a68b0d047344fdb97c510b64a7a3e2d
		NotionStartPageID: "2a68b0d047344fdb97c510b64a7a3e2d",
	}
	bookReactNative = &Book{
		Title:          "React Native",
		TitleLong:      "Essential React Native",
		Dir:            "reactnative",
		CoverImageName: "Ruby.png", // TODO: need the cover
		// https://www.notion.so/kjkpublic/Essential-React-Native-c7980909d5144eb5aee8b28bbe60ec9b
		NotionStartPageID: "c7980909d5144eb5aee8b28bbe60ec9b",
	}
	bookRuby = &Book{
		Title:          "Ruby",
		TitleLong:      "Essential Ruby",
		Dir:            "ruby",
		CoverImageName: "Ruby.png",
		// https://www.notion.so/kjkpublic/Essential-Ruby-51f7633cdf1f4ab1a7778e8095f599bd
		NotionStartPageID: "51f7633cdf1f4ab1a7778e8095f599bd",
	}
	bookRubyOnRails = &Book{
		Title:          "Ruby On Rails",
		TitleLong:      "Essential Ruby On Rails",
		Dir:            "rubyonrails",
		CoverImageName: "RubyOnRails.png",
		// https://www.notion.so/kjkpublic/Essential-Ruby-On-Rails-80d02f56455d4162a91223e5fc1341e0
		NotionStartPageID: "80d02f56455d4162a91223e5fc1341e0",
	}
	bookSwift = &Book{
		Title:          "Swift",
		TitleLong:      "Essential Swift",
		Dir:            "swift",
		CoverImageName: "Swift.png",
		// https://www.notion.so/kjkpublic/Essential-Swift-e76d42906b0e493291a60bbd351f3b6b
		NotionStartPageID: "e76d42906b0e493291a60bbd351f3b6b",
	}
	bookTypeScript = &Book{
		Title:          "TypeScript",
		TitleLong:      "Essential TypeScript",
		Dir:            "typescript",
		CoverImageName: "TypeScript.png",
		// https://www.notion.so/kjkpublic/Essential-TypeScript-9f3a0df9855747b1ab85b76637971d62
		NotionStartPageID: "9f3a0df9855747b1ab85b76637971d62",
	}
	bookBatch = &Book{
		Title:          "Batch",
		TitleLong:      "Essential Batch",
		Dir:            "batch",
		CoverImageName: "Batch.png",
		// https://www.notion.so/kjkpublic/Essential-Batch-cmd-exe-ea84bde7ed4e4353bdc6ae44125abc08
		NotionStartPageID: "ea84bde7ed4e4353bdc6ae44125abc08",
	}
)

var (
	booksMain = []*Book{
		bookGo,
		bookCpp,
		bookJavaScript,
		bookCSS,
		bookHTML,
		bookHTMLCanvas,
		bookJava,
		bookKotlin,
		bookCsharp,
		bookPython,
		bookPostgresql,
		bookMysql,
		bookIOS,
		bookAndroid,
		bookBash,
		bookPowershell,
		bookBatch,
		bookGit,
		bookPHP,
		bookRuby,
		bookNETFramework,
		bookNode,
		bookDart,
		bookTypeScript,
		bookSwift,
	}
	booksUnpublished = []*Book{
		bookAlgorithm,
		bookC,
		bookObjectiveC,
		bookReact,
		bookReactNative,
		bookRubyOnRails,
		bookSql,
	}
	allBooks = append(booksMain, booksUnpublished...)
)
