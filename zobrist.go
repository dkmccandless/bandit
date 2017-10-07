package main

type Zobrist uint64

// Cryptographically random bistrings for Zobrist hashing. Do not mutate.
var (
	pieceZobrist = [2][7][64]Zobrist{
		{{},
			{0x5d37bab864876915, 0xcfb096455696e570, 0x895eb16a2ac5d30b, 0x42cc4ef65a8e22f7, 0x8bc361bbe045ec84, 0x83a91527d396e545, 0xfaec68704683263, 0xa922bc2279856cc0, 0x1c7fcc2672040137, 0x94567069ff669ab0, 0x183f3b7d5f35398f, 0xdd699383998e246f, 0xbc26d7db84f99077, 0xb9e865d572394411, 0x4aba3671dc370e1f, 0xe7d78ca86cc5be14, 0x7f473734617fb5fa, 0xe0cc71a4f812f025, 0x84c47dc7727a8c25, 0x9c7ac67bdc747ab1, 0x6397dff385207a4c, 0x91ca32e47866ca61, 0x592e16f579e1466, 0x98d0b0c746f31e27, 0x5fac7fe23631a654, 0xe6189a74432cf5ac, 0x2907a5f73ee6dea9, 0x5c8c124339e6195d, 0xca02ea959ad391f, 0xa78d6d1b81ad38c5, 0xbdc0a347a088561e, 0x867a4e6bd80b548b, 0x5f6dc9f99583d4a5, 0xbd7475554b79f63b, 0xfe2dd6eea9fb7bd4, 0x810244537d12da97, 0x27129f70e3fd3dd6, 0xa30dcb27a7e9074e, 0x82b2b2b4dba99576, 0x3d87ea1a3cefee9e, 0x14173b3ca2f9fe2d, 0x5007f805aee40889, 0x2ba3c608eac311c6, 0xfc32be24e85d1af0, 0x74513161cbf552b3, 0x674bbe1d8d9a65f0, 0x88a8856732794bc, 0xfdf1f9f9df4646a1, 0x714a1e0f2e017408, 0xedf2061b9183ac63, 0xafefe8f5ba50d5d9, 0x8685945b32b9aba9, 0x27a81b84233e7cab, 0x424d9c312fd8b96c, 0x5303b27de98e0778, 0x51a0e980ccc4e23b, 0x8d57f0d84551e068, 0x2d2a351027f6d4af, 0xd7cc526abf7e0780, 0x388470dfb959cc37, 0xefbe799b6a796e49, 0x81929d9821891b24, 0xf519b44e32585825, 0xc428f5f283e8d85f},
			{0x4e975fe2fb39c5bb, 0x72719ae797c61b1b, 0x523409ca4a2077ea, 0x75f9750290d81684, 0x5ec6f93bfe878754, 0xf99dfec24d8280f8, 0x707bb56171b66ed4, 0x520f0e19b3c8bb89, 0x6871c180c6407e77, 0x78079a1c89cf6b56, 0xa14b5f9af2a6b2c5, 0x12a021bcee87ac1f, 0x585f2e22139073f4, 0x9ba046c3ba8a6f39, 0xc3d7fccdb16d31a6, 0x9f5e77b30872900f, 0x14fa3f3f82bd7fe3, 0x2eb182d2a4690bd5, 0xb251161277b53a1f, 0xdd107ea5f9e4744e, 0x1a954622b1365236, 0xb8a4f99f239a4d19, 0xb2ffa4380f3ac0fc, 0x76ea8c08ea7b6485, 0xfda0ccb40a062a8f, 0xdffd9dbb5f2b56f3, 0x8d3c63bbfe9ec6cc, 0xe8efd5464f3dbf90, 0xb12cb0b5ccfd95a2, 0x900ef66698d552e9, 0x7495e82fd344ea34, 0x703ff3ed5541d349, 0xb8e6030c60ebfdc3, 0xedf5e3835e31d6c0, 0x6b6fb13cb59daddd, 0x65e0f02d47ea135e, 0x19574e608bda0950, 0xd3703a7bad473931, 0x9731b6e975f1581e, 0x18c0d69bf5be2d29, 0x1fcfc6b95deb1ade, 0x9d141e9d376f8c3a, 0x4372cdd4a304a30e, 0xd5adb201c179a804, 0xef6878a62436d0f1, 0xea3983155c9d52fd, 0xd6bb13cacd01c40c, 0x171f04afd67bb48a, 0x6aca4f76f4e6de1d, 0x990a68ac71427d34, 0x864c48e724479195, 0xd0c7a0c526d13c4d, 0x3a9ed61ce6665565, 0x3817b4a6a4274b45, 0x8698676c23e5e561, 0x3ba75fd36bbc90a7, 0xb525ab5313d5badc, 0x66b83570c2558419, 0x33fb7003ce58eaae, 0x268fa3d0394ed00a, 0x9c5dcee073149af3, 0xf4362393fe487367, 0xce5cf6ae43ed692f, 0x56cf81a1c7b35002},
			{0xb71a2f46a3a71ae9, 0x7ec9a8ca64cd7e97, 0x724da930c0a5d74f, 0x404834861f55a59f, 0x72269eb4135397aa, 0xdb6441e333dbf0b3, 0x93c4a7c1730e512f, 0xe769744dfcb89ecc, 0x9ffc5be25f831e6e, 0x974683898ca312e8, 0x20b535cfbf12d768, 0x8a77e39dcf0ada9f, 0x5d04d3277a450241, 0xeae6c118502f4b89, 0xf1c07f8170266b75, 0x6dd001087413d214, 0x9d7ff220ecea61c9, 0x353d06815b4c8161, 0x1465b91c86de7ef6, 0x44a123d1e43346e5, 0xe29b0882421c4193, 0x2ca6b0720ca63d9e, 0x3ee10e6c5a39857b, 0x4fb12e876f466b1a, 0x21656770a5a01bbe, 0xabfaf289d5a278eb, 0x7cfd92dd25a2d536, 0x5d8b33a27071a7d4, 0xe9f1aa8e48f74d93, 0x897795d61ff72564, 0x8175a80ded6dfe61, 0xd7354dc06410bb85, 0x8877cf2dea342a6a, 0x4fd808281b1f8ed5, 0x89678357e18985d9, 0x555697123d87a723, 0x1d5445e346527e3b, 0x29e4da9fb6c2406b, 0x89eb776c9f3c37c4, 0x1de605717bd383f6, 0x631bfa9a1d11633f, 0x8b9cfe53a7b3e532, 0xfa73f76111607a63, 0xcf028e1b12fe07f9, 0x8ba31ca0306bcc8c, 0xe7ab145e55fd4ffb, 0x461c3678f2fa45ab, 0x9fb3f1dec8898b61, 0x48705f07a520a4f5, 0xa520eede3a631d9f, 0xcea7941834f9cb73, 0xfd4d95110f646f3f, 0xd9f92c8e4f67ffeb, 0x8c22a7dbe926250d, 0xe76511b72076eca1, 0xbe5432ff9b5bd31e, 0xfee36354221b2593, 0x7598a9ba00ce6704, 0xd4496e431a6c1fb5, 0xa3fd6ce86f696763, 0xd2c7b7dd10eaef4d, 0xeb9b0befbe9b174f, 0x4452ac7084d1cef3, 0x820b5b998499b446},
			{0x5bf1052c2e27c5bf, 0x282f1ff7c4bb2e65, 0xa18f4dbaf9ec868f, 0x86e609535443772d, 0x883b2623535b8be6, 0x73ecc60b8d4b2aa, 0xaf08196d5cb1fac2, 0x8597b6974868e9ea, 0x54191379aa21b981, 0x32aa38b0175fb185, 0xebccc6bf8ae5bc3c, 0xa0a2d1b14eaf65d2, 0xc5c4d023fbeef464, 0xb199cea86ca026e2, 0xc6a65d548849a72d, 0x3b6b9979ec78fd6c, 0x39828b59aa7d7351, 0x4407f26d719fc7f9, 0x5a31098dd11a3f28, 0x63f2ce92382c95ef, 0x63417294c36aa0d0, 0x614dc2ad25bfa225, 0xc19230a4334262c2, 0xd2f2a53a30dd3a86, 0xf34ffb046c2cedf5, 0x50ce5666f83ed3cd, 0x6291fae9eb1fe5ad, 0x19f8ea1a2574aaca, 0xddc8ef7041edf17d, 0xe9ee3d4794470ad, 0x93b21ffeb8e24858, 0x977108ac18c28ab5, 0x64925868bbe8ec96, 0xc91e4e738d38874a, 0x54a155d705557624, 0x3a1f34d2621ac977, 0x47c38b6d58c72dfe, 0xd57b43bc72643de5, 0x5c329d2a22849cd, 0x4c8e80888eaceaa8, 0x8fe4545122ed70bf, 0xaf99a49761cf95b9, 0xb3f61d6dd13be951, 0xb0601dac3a1b1266, 0xefb4931a4f5cc0c1, 0xce05062de2139db0, 0x7ea9876933262c0a, 0x63ffe47ba896df08, 0x5c79b74c570e60ef, 0xd081ce278bc5936a, 0x3d83de66e0cd2404, 0xbf307d0b4ba69b6a, 0x47ca930976ee30e5, 0xf20cc00eac934869, 0x8d8fed78900fed54, 0xfe857bce9098af15, 0x52ba53a255d0b3a8, 0xdba836685d7cf5e4, 0x2b707632552a3c2b, 0xf637a4354d5f9609, 0x6e1392ef4db4da2b, 0x89b8e6e4ccbcec85, 0x416b7092e9e9e60e, 0xdb969b997ebc5da9},
			{0xf9d1e673100e1673, 0x9b5323916905a995, 0x13405df8c82efa33, 0xb7abe4d64eaf2436, 0x9c20490a7e63e06f, 0xfd5d510c9c501b59, 0x94ef7eb967abc618, 0x7d057fc316e0eaf, 0xbdb96464060f5fc2, 0x5aac47592c7e4539, 0x96ee5bbcfe7d8939, 0x724f1127e8dbb093, 0x6def3583b00a453d, 0x8eaa2099048ab9fd, 0xd42f673638879c0, 0xcce678e8f0e9c01d, 0x9fb6cf0ced0d9842, 0x956b7fb5b89f41b4, 0xccbfe596e1be1687, 0x761b33377ec36e8d, 0x665074e1278955fd, 0x218005d9a2ff49d9, 0xe01f120e7a778e7f, 0xbbe9d3970f9ba8ea, 0x813c4fee09276e33, 0x45a1fad65a54d66d, 0x2f313569a67923a9, 0x77da58b399961ffe, 0x627b89875c93433d, 0xb1a2c048ca5657cd, 0x4994ab26c5fefe34, 0xc1eadda45ba387, 0x53411e886163c126, 0x91492c5e9474c4a7, 0xcec98a240a2e4051, 0x83cd2b4b1b7cf636, 0x1ac1ff059d9aab3f, 0x45a91b1b0b0779b8, 0x63f98b67538bcd33, 0xe5c5396acab3238, 0x41cc3ea673f9a27d, 0x3e0a04251964b1c0, 0x99c49fb915666561, 0xa4393dc96085e60e, 0x5cff3ae3bb8aa38a, 0x6b81654934e4e4, 0xde9a01abaeb0e069, 0x409a061c849a9a43, 0x2b7d5fd530c0ac3a, 0xb68344ff232abfbd, 0x50eb46bc0429426b, 0x9e1a42abb222ccc9, 0xda2717844b0a969a, 0x31bda102c22a3d78, 0x1697be0f0efaf97a, 0x86ad1280c100a4b0, 0x769bb046749c6681, 0x7ea86462bd84702d, 0x7bf5f5d303d97065, 0xaa7e13cf701aaae3, 0x3b36bb35efe7def1, 0xb3e4df07e5b114fc, 0x4bb63070995dd9c7, 0xd3f5ce5d832d1bc},
			{0xe1ff598e08c7306e, 0x82d9f4a18623e409, 0x675bffc8a68e85a, 0x629df1c9c1da03ff, 0xa3641340827a176d, 0xe0f90d8c568b9f20, 0x9ffbd13f89ec7449, 0x66399e7b2cebadfe, 0x542bf0666148a728, 0xb1e2c6882d466dc6, 0x5ab5f84cd59a22d0, 0xeeaab70cd83377e, 0x780390ae9ece5fd0, 0x36952c03dd74207f, 0x846c6fc1fcce3c5f, 0xf0fbee5526fa3e95, 0x25181f8c1c0f2080, 0x72a205030d7eaca2, 0x611d9d6bbc57ae08, 0xbee0b4f4e5c247ca, 0x2027ac3fdbd10e95, 0xe348b6e5d0198d3a, 0x449e4143b20a48a7, 0x64c63782a9318d92, 0x13fb4503cab1e32e, 0x83c6bfe41ccabac1, 0xedc65db997205f1e, 0xb5c571cc2ee4928b, 0xc8fe6798cade73df, 0xdc30fb8024083418, 0x9df5b7d04ba9c8eb, 0x749d0d82cf40ec47, 0xde2c7db332f4123, 0x5ea8f987fde524fa, 0x57d453af6413e978, 0x51444b41383b1d18, 0x41133394a6cc51d3, 0x475d73865c09fbff, 0x62ff47fedd84aa83, 0xf8718a37d7639ed2, 0x2b8b4b1e39fcee51, 0xbb54c31dfbb4de4a, 0xbde9d960146be5e3, 0x1cb3765f164fac8f, 0xb02a2c275d7adc2c, 0xc939ffb6ea15e3bc, 0xa618642f4c196965, 0x854231c60c16cb76, 0xdef0518e68a4fb1b, 0x23af36a0bacd57bb, 0xd96ea7ef18ef7ec7, 0x4ed4f8211e6daafb, 0x197c54c13b9fed09, 0x909e3b7ab1fef800, 0x9248479bbc483d92, 0x5781b9d1e410c1bd, 0x472257a69ded5930, 0xbc58dd07b9013de6, 0xb538bcced516640f, 0x2149623d9a244333, 0xbeeb490c265be560, 0x8e9b9f5ecbb711b0, 0x2029cf0eec14cbda, 0x6ad2d266b78087cb},
		},
		{{},
			{0xfa27a23739390ba1, 0x7935af0141a8d79d, 0x517a35280cef0486, 0x7f11603faedc0ec4, 0x9a4564358fd6319f, 0x707a26cb88014660, 0x18170d261e9a2301, 0x1b4e664ae54bbefa, 0x2db0f8cc86129d58, 0x847c992f1f1bc21b, 0x1737749bb4f42380, 0x3c20384bd9a6b247, 0x2ae10d65e83f6ef9, 0xc8b9db7b84d78981, 0x25b855a28ca17f5, 0x86209419bff92503, 0x5ea1b4168b468170, 0x1a3d3f3811746b3c, 0xa9052a50f62055fe, 0x2efd9f14bf8bb284, 0x732a5939ee1309c7, 0xc76dee9cb8e3e721, 0xeb2417b1f0b82da8, 0x70d663171ec5df27, 0xbf2fc3e0f86307fc, 0xe799b37bae22b76, 0xf711c097658d0648, 0xed3e80aa704b2a4b, 0x1e3b76235bdcc039, 0xf5e6bc475b55e550, 0x7d7c5df8a7a82e85, 0xbae1070698318f9c, 0x5a7092332487746c, 0x7a8618b9b02ef242, 0xf76178965dd97c07, 0xa2dd68ebabe9e038, 0xac22dfec43eedc49, 0x2da0dfc3e78af7a7, 0xf7474104876c6209, 0x71e537bb4ec726d4, 0x3255fec4cbd36c89, 0xa337ac0684529341, 0x8978fc949c644a79, 0xb5072c9fcbaec4b9, 0xb0de160c3e36d6a8, 0x34d33c89685ced3a, 0x4fad9db8c82710ef, 0xd5a851d6bb9b5bbf, 0xbd2b6b11c3969aee, 0x2597b359558a17c8, 0xe11a24f01fe26060, 0x1d4876a695eaf6f4, 0xa8a498b4c7aff41, 0x5ff8c16dd33dc3d1, 0xc17b8b2fc3fc5bfa, 0x2e00d46c3655839d, 0x3b2b6e8b55969cdc, 0x75fe90b17916f7f2, 0xbed351e7d683494f, 0xec23fbdeab19a3a5, 0x2274f4f68c10cfb9, 0x9d7c4ea8f2bb0568, 0x5e3e75846e2df480, 0xf7ce62a339720fe2},
			{0xa069a08c9626ff11, 0x3c809877a4ef7684, 0xd278d4f370d26f83, 0x8926bb950be8c763, 0xefee22d30e746347, 0xfc7589eedd514025, 0x2953c96385dbe2e0, 0x559d6808e38b2fe6, 0xd3108443d6ea99e9, 0xdb80ec7d100a4254, 0x9562a7c457059e11, 0x1c1ff9e27efcce59, 0x564ed5ed7e23b1a5, 0x4cab1245db1d0163, 0x3b12495e5e8743fd, 0x10240f2d29ce4ff, 0x7ab12737e605df5a, 0x8971ad25327a0fc, 0x6a8396f6a89be729, 0x7506aaad84e65e26, 0x3878dab77a932991, 0x5f8e2fde78163e2e, 0xc3d7a4a54b5fbda5, 0xe6f6cced6dfcc7f8, 0x32a82aa9b8891511, 0xec6090df5d5db7cb, 0xf42d5fde975fda98, 0x960f4d48cf51de99, 0xdf377690b2edb059, 0xcadf54bbaa87e7c3, 0xdff95f89c46d3b3, 0xe84caaea10ac6ecc, 0x6c96d7614930f6da, 0x2e87cb9288f4fe91, 0x2a2cf09567c0596a, 0x6270b6e277159055, 0x97b0b7af60fe7ee3, 0xeb3d8484027d987b, 0x8ced54632d9e0e8e, 0xf208d7b6fa7f49e7, 0x4bccb81dc04142aa, 0x4a1c12b12b3e5a45, 0xa5f6770a7d14894c, 0x5c7a5b09b4ff6f, 0xc2d3d95748ca28fc, 0x537c7987c853f178, 0x93c9263b356863f2, 0x6c467afe0ae92fbb, 0x973947db1ccd7cbf, 0xe36fd26038ec0f55, 0xac67100e0509120e, 0x8add30312b29b3d2, 0xcb7af14eac4e5f70, 0xf3b6db0e2c901651, 0x65fa9c11a6dd4068, 0x54602c01c951cbc2, 0x86daf268014f26b1, 0x5f6e84c186fa3f9b, 0x1a4f4410cccfd7d, 0xd777e7dba91683ec, 0xbf8000d14fdc8afa, 0xf8a49941cdfc34ee, 0xfb5a5c2a971a5361, 0xe6abf9eb2e3e061d},
			{0x97e046db9b561963, 0x858f61cf2e3a4bd4, 0x805c7b980b2b874e, 0xe6e2499bbb60e654, 0x2793ebf53d2a9571, 0xd3f14a5d0ab61cfd, 0xf3354f29b605f1b3, 0x66781e413c9094c, 0x7680798340f996bc, 0xfd317a7d7412e531, 0xd0635aa7ce66f88f, 0xd1e7ef8e336c3fb1, 0x5294e5fac67ef8d3, 0xa3c344f77a493c4d, 0xbd043b272f3a7168, 0x93f1f0e2fcb66692, 0x352b8739f268cccf, 0xcd9d2221d4755e8, 0xaa4411a1bee76421, 0xfa3b1de61d02bf58, 0x6b3be6b44adb74e1, 0xb52b6098d460d8f5, 0x8b55c7fbff861802, 0x4e616774d762705c, 0x136bf356358bdeda, 0xf64ad8994dae3d08, 0x42bee6e07ad85373, 0x637dba8dda368f57, 0xb9586a5a0f9af76d, 0x7bb787b952f0c291, 0xf682faa82dca518f, 0x49a8aa0ffa0555fd, 0x2540eb0195324a96, 0x9aa903c250ab8645, 0x501ec96458b90907, 0xdcbf407ac3a6232b, 0xc418cee8940f959e, 0xfef408ae16feadcc, 0x3075e7586684285d, 0x9752b93378e95168, 0xfc3ad29aacfa528c, 0x3c7db4525e31c3b4, 0x3a97ac0869280cac, 0x6eb65aab7acab902, 0xa5ef17b4f71494b2, 0xa9d65d854b4ccaee, 0xf2eef338e5267e55, 0xe44e53be3b2cfe4e, 0x603791cc7deeb2ec, 0x1815d27b21301fe9, 0xb3678ad07d2d8d69, 0xb11aede194544399, 0x72370bb579650c0d, 0xbe38dcb1b26a94d2, 0x5044e02eef34a87d, 0x1f1d170df9941339, 0x2ca3d46b6ca1220c, 0xf526813417e9a681, 0xd750799beba69874, 0x355bc09efcba36f1, 0xe799449b557fc739, 0x9f66c8f0d26febfe, 0x1530e1e3c7508552, 0xb733aee5dd2da04d},
			{0x3daef1910b525ab5, 0xcd871a4e90b621d7, 0xefa72f914a12f2e3, 0xa7cc53fe75ce54ba, 0x789dc6fe32fc2a75, 0xadeb73f738dd006e, 0x501f62fc8c985b89, 0xe9c90d372464b117, 0x79d96ce631a91816, 0x5915447883090309, 0x81113aff34eb2324, 0x8a12129ae87db2df, 0xe8d0a80cf29f095c, 0xe226852e266539b, 0xcece5968d16cbf18, 0x4f81b5325ad54b2, 0xb9f74d81dde87524, 0xa8918f37ce1b45d5, 0x3ede8aa70d341119, 0x71c9b8f44622a606, 0x71738721f9ccac21, 0x420f9736268c3b07, 0x9309ad36bd91be57, 0xbc855660b526448b, 0x6c20ff9a27de4c6b, 0xb28af77a19a9b7e, 0x8533a6626fa552d4, 0x83a83981f9ad0c7e, 0x1cb4504276ebc3a3, 0xad99e14c9afaafb, 0x642acc8aa6a8d739, 0x9ed239e77164b165, 0x93fdba31b8eb98b1, 0xdb2669f645668cf7, 0xcd2f9a974dd463b3, 0xc343624d25d31bea, 0x4bdd54e2df5b22b2, 0x1dbe5839dc704c3c, 0x5f098c47a2f19d, 0xc44516cec23bba98, 0xb67730bb5dda6209, 0x888ccfa7ac2db338, 0x46ddfce9d6f1feb0, 0x9a8efde6087a24f9, 0x8427142ddc5ba749, 0x29585e47e5882ee2, 0xa5b0b7626a9f0f95, 0x56c10b32869cca33, 0x538b0765fbb533be, 0xa5bb4bdec92f34b1, 0xa452a43ad58d50dd, 0xc5b4ea4561ab11b5, 0x965ad4134bdcfbc4, 0xbdf85d46e9e189ed, 0xfce72641a756fba5, 0xa84f5cff82717265, 0x628e296c2994c066, 0x3c6ca7472541c093, 0xbbc2229eae8d7426, 0x111fd72dde1fcd97, 0x66b83c29102818ef, 0xe7963e812cf2622d, 0xa1b02183adb0ba09, 0xb7992323cb88e4a1},
			{0xc37ce0eddfed88bc, 0x2c9d824194392111, 0x6555329e27debba0, 0x30886939e627386d, 0x4f3ed0fbed0e7dcf, 0x1d9d7cfe7b714ef0, 0x5617cf6582656d2d, 0x9617e0dc4821db59, 0x2af422c6c2a26272, 0x42d8b71ca342323e, 0x56e719914db583ff, 0x4ebb0c6fce62e2fa, 0x8274fe665f8daa88, 0xe2eaa877498bafdd, 0x33ed1843742eed93, 0xbe2e6317e3b7dc31, 0xca7da3347ef5bc4, 0x2062efab8cab39e3, 0xc3ae6df1bfac3bc8, 0xaa7c6a5a2197b1d6, 0xe339ee17777dab96, 0x6cedc7ddad9ad11d, 0x11cda94105f08a8e, 0x4db79e9d3c11c86, 0x32d5169f5a3e7061, 0x44c0266966b4339a, 0x2a231bb0590bea4d, 0x7d290899ffb172c4, 0xafd5d63fce62cc, 0x69c32aa3cb6dbbc5, 0x63ef69dd16f80aa0, 0xe0cb61b5f208e406, 0x3f72315b3aa1d880, 0x6f2735e751060d6, 0xbcc884a569b027db, 0xd59d773cde0964d, 0x1f9e748d086a527c, 0xe83e75d02f092ea8, 0xeeb1796242ca1607, 0xc84184c24050afe2, 0x41984c2da3e91ba7, 0xaf92ce282d95b71c, 0x75875f692709e960, 0x24f93528b217897a, 0x765d65ebad8440c5, 0x321fe35689669da9, 0xb0008388cfde7e0, 0x4983651c4ffec08, 0xb21639ecaee3b301, 0x112aa2d3cd75f72f, 0xf11a1e38f316c538, 0x987a33b50fdcce47, 0xd19638bd9e3d8001, 0x9b2c4033334ea68c, 0x21ba28661ff9d106, 0x480f1484269096bf, 0x19cff7e6cce49552, 0x710a05d8942d962e, 0x433d496478362e7, 0x10a1ae612d5f7aa0, 0x4f6048051ef3895d, 0xd3e8ceb06023e48e, 0x3f4031761380af98, 0xb34aba522e380409},
			{0x511bbf014c760f6e, 0xde7c8a3d93c6f1b8, 0xbf489dae0b8514e2, 0x7be19fb99f45e7ee, 0x82b67bc4b92183c2, 0x77b09813932ac510, 0xb01f89418c890041, 0x601f7408ed4782a2, 0x65bd90719c4488b9, 0xc0ef22ecaea28d10, 0xaec7b429932ec5be, 0xed88943223cee395, 0xe008d4c7f65f0422, 0xec6846df61a70cce, 0xf2ac027afc1dd59c, 0x74d251f4973d44ae, 0xc03aebc4449621e5, 0xe5be1bf7eed4bccb, 0xe8a1f89aef917aed, 0x4eb8d96743a07210, 0x25172a25348ea370, 0x8967099c13116712, 0xb46b64fb995759b1, 0xfbf17bc07606806d, 0x3de39b9944aff898, 0x4e6268f93f99197, 0x75aa38ca082c08f7, 0x7800d67621edab8, 0x30433ad0f8a5a057, 0xe5610a1131782c2c, 0xdea1c93ff2515fc4, 0xff186ea5ac8e9c37, 0x9e7808df376eaf1c, 0x661fffaa47a49b51, 0x801ee6a6420d016a, 0xd37994be2c72cfae, 0x7609f832359edc66, 0x940cbf314cd6d49c, 0xa98ccebce5d23550, 0xc05f24f1350b29e2, 0x8c26ad3ac5615d03, 0x266ed9622ed14b05, 0x23b5b97f6f4f05cf, 0x9110cf79d6060e0f, 0x60674b2a93e1a938, 0x1f2abd003824bebe, 0xe2d9cf4044641603, 0x99697f7c3055b8bf, 0x32283018e74e03c0, 0x896700ca6190db1, 0xed3a2184fe0d2e01, 0x523d459ec7b76fcd, 0xd761bcd6a70c77e3, 0xe05f0d9fdabc1d54, 0x45ea02d98d99eca0, 0xf24aa11cdef6100, 0x71c82c65947bd15, 0xca17f6a61353a9bd, 0xfa01cb0f2cf7a9d0, 0xa6c513a0b7a10279, 0xf8edf255478086b0, 0xa636cba9cb549ed1, 0x40dc4b071fc94396, 0xd1eb206dc97babc1},
		},
	}

	// Castling rights
	qsCastleZobrist = [2]Zobrist{0xfb29d82b9fddf9ad, 0xc69b325e3b43c1e6}
	ksCastleZobrist = [2]Zobrist{0x96dd7a62f6837f56, 0xb920310a317b3639}

	// The ability of a pawn on the 4th/5th rank of the given file to make an en passant capture.
	canEPCaptureZobrist = [8]Zobrist{0xaff9512fa0250760, 0x877d09dd4d2f622a, 0x4934b5d5668a9505, 0xab131aee3c9bcb00, 0x64351fe47f9c649a, 0xc2362d4d1989b0d0, 0x4fdbba027f1ebc77, 0x2083ed95769430c5}

	blackToMoveZobrist Zobrist = 0x62bf223b8ae0d2e7
)

func (z *Zobrist) xor(x Zobrist) { *z ^= x }

func (z *Zobrist) xorPiece(c Color, p Piece, sq Square) { z.xor(pieceZobrist[c][p][sq]) }

// Zobrist returns a Position's Zobrist bitstring.
func (pos Position) Zobrist() Zobrist {
	var z Zobrist
	for sq := a1; sq <= h8; sq++ {
		if c, p := pos.PieceOn(sq); p != None {
			z.xorPiece(c, p, sq)
		}
	}
	for c, ok := range pos.QSCastle {
		if ok {
			z.xor(qsCastleZobrist[c])
		}
	}
	for c, ok := range pos.KSCastle {
		if ok {
			z.xor(ksCastleZobrist[c])
		}
	}
	if a, b := eligibleEPCapturers(pos); a != 0 {
		z.xor(canEPCaptureZobrist[a.File()])
		if b != 0 {
			z.xor(canEPCaptureZobrist[b.File()])
		}
	}
	if pos.ToMove == Black {
		z.xor(blackToMoveZobrist)
	}
	return z
}
