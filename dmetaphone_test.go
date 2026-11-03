package strsim

import "testing"

func TestDoubleMetaphone(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantPri string
		wantAlt string
	}{
		{name: "Smith", input: "Smith", wantPri: "SM0", wantAlt: "XMT"},
		{name: "Schmidt", input: "Schmidt", wantPri: "XMT", wantAlt: "SMT"},
		{name: "Schneider", input: "Schneider", wantPri: "XNTR", wantAlt: "SNTR"},
		{name: "John", input: "John", wantPri: "JN", wantAlt: "AN"},
		{name: "Robert", input: "Robert", wantPri: "RPRT", wantAlt: "RPRT"},
		{name: "William", input: "William", wantPri: "ALM", wantAlt: "FLM"},
		{name: "empty", input: "", wantPri: "", wantAlt: ""},
		{name: "non_alpha", input: "123!@#", wantPri: "", wantAlt: ""},
		{name: "single_vowel", input: "A", wantPri: "A", wantAlt: "A"},
		{name: "single_consonant", input: "B", wantPri: "P", wantAlt: "P"},
		{name: "Thomas", input: "Thomas", wantPri: "TMS", wantAlt: "TMS"},
		{name: "Thompson", input: "Thompson", wantPri: "TMPS", wantAlt: "TMPS"},
		{name: "Ghost", input: "Ghost", wantPri: "KST", wantAlt: "KST"},
		{name: "Phone", input: "Phone", wantPri: "FN", wantAlt: "FN"},
		{name: "Caesar", input: "Caesar", wantPri: "SSR", wantAlt: "SSR"},
		{name: "Gnome", input: "Gnome", wantPri: "NM", wantAlt: "NM"},
		{name: "Knight", input: "Knight", wantPri: "NT", wantAlt: "NT"},
		{name: "Xerox", input: "Xerox", wantPri: "SRKS", wantAlt: "SRKS"},
		{name: "Zhang", input: "Zhang", wantPri: "JNK", wantAlt: "JNK"},

		// C: CHIA (Italian)
		{name: "Chianti", input: "Chianti", wantPri: "KNT", wantAlt: "KNT"},
		{name: "Chia", input: "Chia", wantPri: "K", wantAlt: "K"},
		// C: CHAE (Greek)
		{name: "Michael", input: "Michael", wantPri: "MKL", wantAlt: "MXL"},
		// C: Greek roots CHARAC, CHARIS, CHOR, CHEM at pos==0
		{name: "Chemistry", input: "Chemistry", wantPri: "KMST", wantAlt: "KMST"},
		{name: "Chorus", input: "Chorus", wantPri: "KRS", wantAlt: "KRS"},
		{name: "Character", input: "Character", wantPri: "KRKT", wantAlt: "KRKT"},
		{name: "Charis", input: "Charis", wantPri: "KRS", wantAlt: "KRS"},
		{name: "Chore", input: "Chore", wantPri: "XR", wantAlt: "XR"},
		// C: CH with Germanic/Greek patterns
		{name: "Orchard", input: "Orchard", wantPri: "ARXR", wantAlt: "ARKR"},
		{name: "Architect", input: "Architect", wantPri: "ARKT", wantAlt: "ARKT"},
		{name: "Orchestra", input: "Orchestra", wantPri: "ARKS", wantAlt: "ARKS"},
		{name: "Loch", input: "Loch", wantPri: "LK", wantAlt: "LK"},
		{name: "Bacher", input: "Bacher", wantPri: "PKR", wantAlt: "PKR"},
		{name: "Wachter", input: "Wachter", wantPri: "AKTR", wantAlt: "FKTR"},
		// C: CH with MC prefix
		{name: "McHugh", input: "McHugh", wantPri: "MK", wantAlt: "MK"},
		// C: CH pos==0, not Greek → X/X
		{name: "Chin", input: "Chin", wantPri: "XN", wantAlt: "XN"},
		{name: "Church", input: "Church", wantPri: "XRX", wantAlt: "XRK"},
		// C: CH pos > 0, not MC → X/K
		{name: "Richard", input: "Richard", wantPri: "RXRT", wantAlt: "RKRT"},
		{name: "Machinery", input: "Machinery", wantPri: "MXNR", wantAlt: "MKNR"},
		// C: CZ (Polish) but not WICZ
		{name: "Czech", input: "Czech", wantPri: "SK", wantAlt: "XK"},
		{name: "Czerny", input: "Czerny", wantPri: "SRN", wantAlt: "XRN"},
		{name: "Czar", input: "Czar", wantPri: "SR", wantAlt: "XR"},
		// C: CZ inside -WICZ is the Polish W rule (TS/FX)
		{name: "Horwicz", input: "Horwicz", wantPri: "HRTS", wantAlt: "HRFX"},
		// C: CIA (Italian)
		{name: "Lucia", input: "Lucia", wantPri: "LS", wantAlt: "LX"},
		{name: "Acacia", input: "Acacia", wantPri: "AKS", wantAlt: "AKX"},
		// C: CC + I/E/H (not HU) - Bellocchio path
		{name: "Bellocchio", input: "Bellocchio", wantPri: "PLX", wantAlt: "PLX"},
		{name: "Bocce", input: "Bocce", wantPri: "PX", wantAlt: "PX"},
		// C: CC + I/E with ACCIO/UCCEE/UCCES
		{name: "Accio", input: "Accio", wantPri: "AKS", wantAlt: "AKS"},
		{name: "Succeed", input: "Succeed", wantPri: "SKST", wantAlt: "SKST"},
		// C: CC without I/E/H → K
		{name: "Soccer", input: "Soccer", wantPri: "SXR", wantAlt: "SXR"},
		// C: CK/CG/CQ
		{name: "Brick", input: "Brick", wantPri: "PRK", wantAlt: "PRK"},
		// C: CI/CE/CY (non-Italian)
		{name: "Cecil", input: "Cecil", wantPri: "SSL", wantAlt: "SSL"},
		{name: "Cello", input: "Cello", wantPri: "SL", wantAlt: "SL"},
		{name: "Cedar", input: "Cedar", wantPri: "STR", wantAlt: "STR"},
		{name: "Cyan", input: "Cyan", wantPri: "SN", wantAlt: "SN"},
		// C: CI/CE/CY (Italian: CIO/CIE/CIA)
		{name: "Facial", input: "Facial", wantPri: "FSL", wantAlt: "FXL"},
		{name: "Species", input: "Species", wantPri: "SPSS", wantAlt: "SPXS"},
		{name: "Spacious", input: "Spacious", wantPri: "SPSS", wantAlt: "SPXS"},
		// C: default C
		{name: "Cagney", input: "Cagney", wantPri: "KKN", wantAlt: "KKN"},

		// D: DG + I/E/Y → J
		{name: "Dodge", input: "Dodge", wantPri: "TJ", wantAlt: "TJ"},
		{name: "Badge", input: "Badge", wantPri: "PJ", wantAlt: "PJ"},
		{name: "Edgar", input: "Edgar", wantPri: "ATKR", wantAlt: "ATKR"},
		{name: "Budget", input: "Budget", wantPri: "PJT", wantAlt: "PJT"},
		// D: DT/DD
		{name: "Cheddar", input: "Cheddar", wantPri: "XTR", wantAlt: "XTR"},

		// GH: pos > 0, not vowel before → K
		{name: "Aghast", input: "Aghast", wantPri: "AKST", wantAlt: "AKST"},
		// GH: pos == 0, next is I → J
		{name: "Ghislaine", input: "Ghislaine", wantPri: "JLN", wantAlt: "JLN"},
		{name: "Ghia", input: "Ghia", wantPri: "J", wantAlt: "J"},
		// GH: pos == 0, next is not I → K
		{name: "Ghat", input: "Ghat", wantPri: "KT", wantAlt: "KT"},
		// GH: Parker's rule (B/H/D before)
		{name: "Bough", input: "Bough", wantPri: "P", wantAlt: "P"},
		{name: "Hugh", input: "Hugh", wantPri: "H", wantAlt: "H"},
		{name: "Doughboy", input: "Doughboy", wantPri: "TP", wantAlt: "TP"},
		{name: "Aughton", input: "Aughton", wantPri: "AKTN", wantAlt: "AKTN"},
		// GH: UGH with C/G/L/R/T before → F
		{name: "Laugh", input: "Laugh", wantPri: "LF", wantAlt: "LF"},
		{name: "Cough", input: "Cough", wantPri: "KF", wantAlt: "KF"},
		{name: "Rough", input: "Rough", wantPri: "RF", wantAlt: "RF"},
		{name: "Tough", input: "Tough", wantPri: "TF", wantAlt: "TF"},
		{name: "Laughter", input: "Laughter", wantPri: "LFTR", wantAlt: "LFTR"},
		// GH: pos > 0, w[pos-1] != 'I' → K
		{name: "Coughlin", input: "Coughlin", wantPri: "KFLN", wantAlt: "KFLN"},

		// G: GN combinations
		{name: "Gnocchi", input: "Gnocchi", wantPri: "NX", wantAlt: "NX"},
		{name: "Agnew", input: "Agnew", wantPri: "AKN", wantAlt: "AKNF"},
		{name: "Segment", input: "Segment", wantPri: "SKMN", wantAlt: "SKMN"},
		{name: "Wagner", input: "Wagner", wantPri: "AKNR", wantAlt: "FKNR"},
		{name: "Signer", input: "Signer", wantPri: "SNR", wantAlt: "SKNR"},
		// G: GLI Italian
		{name: "Tagliamento", input: "Tagliamento", wantPri: "TKLM", wantAlt: "TLMN"},
		// G: initial GY
		{name: "Gypsy", input: "Gypsy", wantPri: "KPS", wantAlt: "JPS"},
		// G: initial GE/GI
		{name: "German", input: "German", wantPri: "KRMN", wantAlt: "JRMN"},
		{name: "Giles", input: "Giles", wantPri: "KLS", wantAlt: "JLS"},
		{name: "Gift", input: "Gift", wantPri: "JFT", wantAlt: "KFT"},
		// G: Parker's rule -GER/-GEL
		{name: "Lager", input: "Lager", wantPri: "LKR", wantAlt: "LJR"},
		{name: "Tiger", input: "Tiger", wantPri: "TJR", wantAlt: "TKR"},
		{name: "Danger", input: "Danger", wantPri: "TNJR", wantAlt: "TNKR"},
		{name: "Ginger", input: "Ginger", wantPri: "KNKR", wantAlt: "JNJR"},
		{name: "Eager", input: "Eager", wantPri: "AKR", wantAlt: "AJR"},
		{name: "Singer", input: "Singer", wantPri: "SNKR", wantAlt: "SNJR"},
		// G: AGGI/OGGI → Italian
		{name: "Aggie", input: "Aggie", wantPri: "AJ", wantAlt: "AK"},
		{name: "Oggi", input: "Oggi", wantPri: "AJ", wantAlt: "AK"},
		// G: Germanic VAN/VON/SCH + GE/GI
		{name: "VanGelder", input: "VanGelder", wantPri: "FNJL", wantAlt: "FNKL"},
		// G: George, -IER suffix
		{name: "George", input: "George", wantPri: "JRJ", wantAlt: "KRK"},
		{name: "Glacier", input: "Glacier", wantPri: "KLS", wantAlt: "KLXR"},
		// G: GG
		{name: "Ogg", input: "Ogg", wantPri: "AK", wantAlt: "AK"},
		{name: "Eggs", input: "Eggs", wantPri: "AKS", wantAlt: "AKS"},

		// J: JOSE/SAN
		{name: "Jose", input: "Jose", wantPri: "HS", wantAlt: "HS"},
		{name: "SanJose", input: "San Jose", wantPri: "SNHS", wantAlt: "SNHS"},
		{name: "Josefa", input: "Josefa", wantPri: "JSF", wantAlt: "HSF"},
		// J: initial J
		{name: "Jan", input: "Jan", wantPri: "JN", wantAlt: "AN"},
		{name: "Javier", input: "Javier", wantPri: "JF", wantAlt: "AFR"},
		// J: vowel + J + A/O
		{name: "Maja", input: "Maja", wantPri: "MJ", wantAlt: "MH"},
		// J: final J
		{name: "Raj", input: "Raj", wantPri: "RJ", wantAlt: "R"},
		{name: "Bajaj", input: "Bajaj", wantPri: "PJJ", wantAlt: "PH"},
		// J: J with adjacent consonants
		{name: "Hajek", input: "Hajek", wantPri: "HJK", wantAlt: "HJK"},
		{name: "Hajj", input: "Hajj", wantPri: "HJ", wantAlt: "HJ"},
		{name: "Julio", input: "Julio", wantPri: "JL", wantAlt: "AL"},

		// L: ILLO/ILLA/ALLE Spanish
		{name: "Cabrillo", input: "Cabrillo", wantPri: "KPRL", wantAlt: "KPR"},
		{name: "Gallegos", input: "Gallegos", wantPri: "KLKS", wantAlt: "KKS"},
		{name: "Vanilla", input: "Vanilla", wantPri: "FNL", wantAlt: "FN"},
		{name: "Valle", input: "Valle", wantPri: "FL", wantAlt: "F"},
		// L: normal LL
		{name: "Gallo", input: "Gallo", wantPri: "KL", wantAlt: "KL"},
		{name: "Allegra", input: "Allegra", wantPri: "ALKR", wantAlt: "AKR"},

		// S: ISL/YSL silent
		{name: "Island", input: "Island", wantPri: "ALNT", wantAlt: "ALNT"},
		{name: "Carlisle", input: "Carlisle", wantPri: "KRLL", wantAlt: "KRLL"},
		{name: "Carlysle", input: "Carlysle", wantPri: "KRLL", wantAlt: "KRLL"},
		// S: SUGAR
		{name: "Sugar", input: "Sugar", wantPri: "XKR", wantAlt: "SKR"},
		// S: SH + Germanic (HEIM/HOEK/HOLM/HOLZ)
		{name: "Shoholm", input: "Shoholm", wantPri: "XHLM", wantAlt: "XHLM"},
		// S: SH normal
		{name: "Ship", input: "Ship", wantPri: "XP", wantAlt: "XP"},
		// S: SIO/SIA
		{name: "Mission", input: "Mission", wantPri: "MSN", wantAlt: "MSN"},
		{name: "Asia", input: "Asia", wantPri: "AS", wantAlt: "AX"},
		{name: "Persian", input: "Persian", wantPri: "PRSN", wantAlt: "PRXN"},
		// S: SIA with slavoGermanic
		{name: "Kasia", input: "Kasia", wantPri: "KS", wantAlt: "KS"},
		// S: initial S+M/N/L/W
		{name: "Smog", input: "Smog", wantPri: "SMK", wantAlt: "XMK"},
		{name: "Snow", input: "Snow", wantPri: "SN", wantAlt: "XNF"},
		{name: "Slow", input: "Slow", wantPri: "SL", wantAlt: "XLF"},
		{name: "Swim", input: "Swim", wantPri: "SM", wantAlt: "XM"},
		// S: SCH + Dutch (OO/ER/EN)
		{name: "School", input: "School", wantPri: "SKL", wantAlt: "SKL"},
		{name: "Schoenberg", input: "Schoenberg", wantPri: "XNPR", wantAlt: "XNPR"},
		{name: "Scherzo", input: "Scherzo", wantPri: "XRS", wantAlt: "SKRS"},
		{name: "Schenectady", input: "Schenectady", wantPri: "XNKT", wantAlt: "SKNK"},
		// S: SCH initial + non-vowel
		{name: "Schnapps", input: "Schnapps", wantPri: "XNPS", wantAlt: "SNPS"},
		{name: "Schiller", input: "Schiller", wantPri: "XLR", wantAlt: "XLR"},
		// S: SC combinations
		{name: "Scion", input: "Scion", wantPri: "SN", wantAlt: "SN"},
		{name: "Scala", input: "Scala", wantPri: "SKL", wantAlt: "SKL"},
		{name: "Scuba", input: "Scuba", wantPri: "SKP", wantAlt: "SKP"},
		// S: French final S
		{name: "Francois", input: "Francois", wantPri: "FRNK", wantAlt: "FRNK"},
		{name: "Courtois", input: "Courtois", wantPri: "KRT", wantAlt: "KRTS"},
		// S: SS
		{name: "Moss", input: "Moss", wantPri: "MS", wantAlt: "MS"},

		// T: TION
		{name: "Nation", input: "Nation", wantPri: "NXN", wantAlt: "NXN"},
		{name: "Fiction", input: "Fiction", wantPri: "FKXN", wantAlt: "FKXN"},
		// T: TIA
		{name: "Militia", input: "Militia", wantPri: "MLX", wantAlt: "MLX"},
		// T: TCH
		{name: "Catch", input: "Catch", wantPri: "KX", wantAlt: "KX"},
		{name: "Watch", input: "Watch", wantPri: "AX", wantAlt: "FX"},
		// T: TH + Greek (OM/AM)
		{name: "Matthew", input: "Matthew", wantPri: "M0", wantAlt: "MTF"},
		// T: TH general
		{name: "Anthem", input: "Anthem", wantPri: "AN0M", wantAlt: "ANTM"},
		{name: "Thyme", input: "Thyme", wantPri: "0M", wantAlt: "TM"},
		// T: TT/TD
		{name: "Butter", input: "Butter", wantPri: "PTR", wantAlt: "PTR"},

		// W: WR → R
		{name: "Wroclaw", input: "Wroclaw", wantPri: "RKL", wantAlt: "RKLF"},
		{name: "Wrong", input: "Wrong", wantPri: "RNK", wantAlt: "RNK"},
		{name: "Wren", input: "Wren", wantPri: "RN", wantAlt: "RN"},
		// W: initial W+vowel
		{name: "Wine", input: "Wine", wantPri: "AN", wantAlt: "FN"},
		{name: "Wawa", input: "Wawa", wantPri: "A", wantAlt: "F"},
		// W: initial WH
		{name: "Whiskey", input: "Whiskey", wantPri: "ASK", wantAlt: "ASK"},
		{name: "What", input: "What", wantPri: "AT", wantAlt: "AT"},
		// W: Polish WICZ/WITZ → TS/FX, -OWSKI → F in alternate
		{name: "Horowitz", input: "Horowitz", wantPri: "HRTS", wantAlt: "HRFX"},
		{name: "Tarnowski", input: "Tarnowski", wantPri: "TRNS", wantAlt: "TRNF"},
		{name: "Lefkowitz", input: "Lefkowitz", wantPri: "LFKT", wantAlt: "LFKF"},
		{name: "Wichita", input: "Wichita", wantPri: "AXT", wantAlt: "FKT"},
		// W: final W after vowel
		{name: "Stew", input: "Stew", wantPri: "ST", wantAlt: "STF"},

		// Z: ZH Chinese
		{name: "Zhao", input: "Zhao", wantPri: "J", wantAlt: "J"},
		// Z: ZZO/ZZI/ZZA
		{name: "Piazza", input: "Piazza", wantPri: "PS", wantAlt: "PTS"},
		{name: "Palazzo", input: "Palazzo", wantPri: "PLS", wantAlt: "PLTS"},
		{name: "Razzi", input: "Razzi", wantPri: "RS", wantAlt: "RTS"},
		// Z: Slavic Z
		{name: "Wozniak", input: "Wozniak", wantPri: "ASNK", wantAlt: "FTSN"},
		{name: "Janusz", input: "Janusz", wantPri: "JNS", wantAlt: "ANX"},

		// French R (IER at end, not slavoGermanic, no ME/MA)
		{name: "Cartier", input: "Cartier", wantPri: "KRT", wantAlt: "KRTR"},
		{name: "Belier", input: "Belier", wantPri: "PL", wantAlt: "PLR"},
		{name: "Papier", input: "Papier", wantPri: "PP", wantAlt: "PPR"},
		// French X (silent at end after EAU/AU/OU)
		{name: "Bordeaux", input: "Bordeaux", wantPri: "PRT", wantAlt: "PRT"},
		{name: "Roux", input: "Roux", wantPri: "R", wantAlt: "R"},
		// initial vowels
		{name: "Eagle", input: "Eagle", wantPri: "AKL", wantAlt: "AKL"},
		{name: "Ice", input: "Ice", wantPri: "AS", wantAlt: "AS"},
		{name: "Open", input: "Open", wantPri: "APN", wantAlt: "APN"},
		{name: "Under", input: "Under", wantPri: "ANTR", wantAlt: "ANTR"},
		// BB, FF, NN, KK, PP, QQ, VV, MM
		{name: "Rabbit", input: "Rabbit", wantPri: "RPT", wantAlt: "RPT"},
		{name: "Staff", input: "Staff", wantPri: "STF", wantAlt: "STF"},
		{name: "Banner", input: "Banner", wantPri: "PNR", wantAlt: "PNR"},
		{name: "Trekking", input: "Trekking", wantPri: "TRKN", wantAlt: "TRKN"},
		{name: "Happy", input: "Happy", wantPri: "HP", wantAlt: "HP"},
		{name: "Quake", input: "Quake", wantPri: "KK", wantAlt: "KK"},
		{name: "Vivid", input: "Vivid", wantPri: "FFT", wantAlt: "FFT"},
		{name: "Hammer", input: "Hammer", wantPri: "HMR", wantAlt: "HMR"},
		// XC/XX
		{name: "Exxon", input: "Exxon", wantPri: "AKSN", wantAlt: "AKSN"},
		// PH
		{name: "Philly", input: "Philly", wantPri: "FL", wantAlt: "FL"},
		// H: first letter before a vowel
		{name: "Help", input: "Help", wantPri: "HLP", wantAlt: "HLP"},
		// H: between two vowels
		{name: "Ohio", input: "Ohio", wantPri: "AH", wantAlt: "AH"},
		// Width (has W → slavoGermanic → affects handlers)
		{name: "Width", input: "Width", wantPri: "AT", wantAlt: "FT"},
		// initial silent GN
		{name: "Gnosis", input: "Gnosis", wantPri: "NSS", wantAlt: "NSS"},
		// Witz (slavoGermanic)
		{name: "Witz", input: "Witz", wantPri: "ATS", wantAlt: "FTS"},

		// G: Parker's rule + RGY/OGY exceptions
		{name: "Energy", input: "Energy", wantPri: "ANRJ", wantAlt: "ANRK"},
		{name: "Zoology", input: "Zoology", wantPri: "SLJ", wantAlt: "SLK"},

		// W: final W after a vowel → F in alternate ("Arnow" matches "Arnoff")
		{name: "Elbow", input: "Elbow", wantPri: "ALP", wantAlt: "ALPF"},
		{name: "Bow", input: "Bow", wantPri: "P", wantAlt: "PF"},

		// Reference values, checked against Apache Commons Codec.
		{name: "Dumb", input: "Dumb", wantPri: "TM", wantAlt: "TM"},
		{name: "Campbell", input: "Campbell", wantPri: "KMPL", wantAlt: "KMPL"},
		{name: "Psychology", input: "Psychology", wantPri: "SXLJ", wantAlt: "SKLK"},
		{name: "Yankelovich", input: "Yankelovich", wantPri: "ANKL", wantAlt: "ANKL"},
		{name: "Jankelowicz", input: "Jankelowicz", wantPri: "JNKL", wantAlt: "ANKL"},
		{name: "Filipowicz", input: "Filipowicz", wantPri: "FLPT", wantAlt: "FLPF"},
		{name: "Arnow", input: "Arnow", wantPri: "ARN", wantAlt: "ARNF"},
		{name: "Rogier", input: "Rogier", wantPri: "RJ", wantAlt: "RJR"},
		{name: "Breaux", input: "Breaux", wantPri: "PR", wantAlt: "PR"},
		{name: "Wasserman", input: "Wasserman", wantPri: "ASRM", wantAlt: "FSRM"},
		{name: "SanJacinto", input: "San Jacinto", wantPri: "SNHS", wantAlt: "SNHS"},
		{name: "MacGregor", input: "Mac Gregor", wantPri: "MKRK", wantAlt: "MKRK"},
		// Codes never contain spaces, even where the reference adds nothing.
		{name: "Maharaj", input: "Maharaj", wantPri: "MHRJ", wantAlt: "MHR"},
		{name: "OBrien", input: "O'Brien", wantPri: "APRN", wantAlt: "APRN"},
		{name: "non_latin", input: "東京", wantPri: "", wantAlt: ""},

		// C: " C"/" Q"/" G" (space-separated cluster)
		{name: "MacQueen", input: "Mac Queen", wantPri: "MKN", wantAlt: "MKN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pri, alt := DoubleMetaphone(tt.input)
			if pri != tt.wantPri || alt != tt.wantAlt {
				t.Errorf("DoubleMetaphone(%q) = (%q, %q), want (%q, %q)",
					tt.input, pri, alt, tt.wantPri, tt.wantAlt)
			}
		})
	}
}

func TestDoubleMetaphoneMatch(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		want bool
	}{
		{name: "Smith_Schmidt", a: "Smith", b: "Schmidt", want: true},
		{name: "Smith_Jones", a: "Smith", b: "Jones", want: false},
		{name: "empty_empty", a: "", b: "", want: true},
		{name: "empty_nonempty", a: "", b: "Smith", want: false},
		{name: "Knight_Night", a: "Knight", b: "Night", want: true},
		{name: "Robert_Robert", a: "Robert", b: "Robert", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DoubleMetaphoneMatch(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("DoubleMetaphoneMatch(%q, %q) = %v, want %v",
					tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestDoubleMetaphoneEncoder_MaxLength(t *testing.T) {
	e := &DoubleMetaphoneEncoder{MaxLength: 6}
	pri, alt := e.Encode("Schneider")
	if len(pri) > 6 || len(alt) > 6 {
		t.Errorf("codes exceed MaxLength 6: pri=%q (%d), alt=%q (%d)",
			pri, len(pri), alt, len(alt))
	}
}

func TestDoubleMetaphoneEncoder_Interface(t *testing.T) {
	var _ DualEncoder = NewDoubleMetaphone()
}

func BenchmarkDoubleMetaphone(b *testing.B) {
	e := NewDoubleMetaphone()

	b.ResetTimer()
	for range b.N {
		e.Encode("Schneider")
	}
}

func BenchmarkDoubleMetaphoneMatch(b *testing.B) {
	e := NewDoubleMetaphone()

	b.ResetTimer()
	for range b.N {
		e.Match("Smith", "Schmidt")
	}
}
