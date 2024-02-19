package config

import (
	"encoding/json"
	"log"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"

	"github.com/CubeWhyMC/NoDelay-Proxy-Server/common/set"
	"github.com/CubeWhyMC/NoDelay-Proxy-Server/version"

	"github.com/fatih/color"
	"github.com/zhangyunhao116/fastrand"
)

const DefaultMotd = `data:image/png;base64,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`

var (
	Config     configMain
	reloadLock sync.Mutex
)

func LoadConfig() {
	configFile, err := os.ReadFile("NoDelay.json")
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Configuration file is not exists. Generating a new one...")
			generateDefaultConfig()
			goto success
		} else {
			log.Panic(color.HiRedString("Unexpected error when loading config: %s", err.Error()))
		}
	}

	err = json.Unmarshal(configFile, &Config)
	if err != nil {
		log.Panic(color.HiRedString("Config format error: %s", err.Error()))
	}

success:
	LoadLists(false)
	log.Println(color.HiYellowString("Successfully loaded config from file."))
}

func generateDefaultConfig() {
	file, err := os.Create("NoDelay.json")
	if err != nil {
		log.Panic("Failed to create configuration file:", err.Error())
	}
	Config = configMain{
		Services: []*ConfigProxyService{
			{
				Name:          "凛加速-IPLC",
				TargetAddress: "mc.hypixel.net",
				TargetPort:    25565,
				Listen:        25565,
				Flow:          "auto",
				Minecraft: minecraft{
					EnableHostnameRewrite: true,
					OnlineCount: onlineCount{
						Max:            10,
						Online:         -1,
						EnableMaxLimit: true,
					},
					MotdFavicon:     "{DEFAULT_MOTD}",
					MotdDescription: "§d凛加速-Hyp加速ip§e 爱来自林冽~ 加入Q群439151992了解更多资讯！",
				},
			},
		},
		PrivateConfig: &Something{
			ListAPI:        "http://verify.osunion.top/IsWhiteList.php",
			Header:         "凛冬网络 | 凛加速",
			ContactName:    "官方QQ售后群",
			ContactLink:    "439151992",
		},
		Lists: map[string]set.StringSet{
			//"test": {"foo", "bar"},
		},
	}
	newConfig, _ := json.MarshalIndent(Config, "", "    ")
	_, err = file.WriteString(strings.ReplaceAll(string(newConfig), "\n", "\r\n"))
	file.Close()
	if err != nil {
		log.Panic("Failed to save configuration file:", err.Error())
	}
}

func LoadLists(isReload bool) bool {
	reloadLock.Lock()
	defer reloadLock.Unlock()
	var config configMain
	if isReload {
		configFile, err := os.ReadFile("NoDelay.json")
		if err != nil {
			if os.IsNotExist(err) {
				log.Println(color.HiRedString("Fail to reload : Configuration file is not exists."))
			} else {
				log.Println(color.HiRedString("Unexpected error when reloading config: %s", err.Error()))
			}
			return false
		}

		err = json.Unmarshal(configFile, &config)
		if err != nil {
			log.Println(color.HiRedString("Fail to reload : Config format error: %s", err.Error()))
			return false
		}
	} else {
		config = Config
	}

	for _, s := range config.Services {
		if s.Minecraft.MotdFavicon == "{DEFAULT_MOTD}" {
			s.Minecraft.MotdFavicon = DefaultMotd
		}
		s.Minecraft.MotdDescription = strings.NewReplacer(
			"{INFO}", "NoDelay "+version.Version,
			"{NAME}", s.Name,
			"{HOST}", s.TargetAddress,
			"{PORT}", strconv.Itoa(int(s.TargetPort)),
		).Replace(s.Minecraft.MotdDescription)

		if samples := s.Minecraft.OnlineCount.Sample; samples != nil {
			var convertedSamples []Sample
			switch samples := samples.(type) {
			case map[string]any:
				convertedSamples = make([]Sample, 0, len(samples))
				for uuid, name := range samples {
					convertedSamples = append(convertedSamples, Sample{
						Name: name.(string),
						ID:   uuid,
					})
				}

			case []any:
				convertedSamples = make([]Sample, 0, len(samples))
				var u [16]byte
				var dst [36]byte
				for i, sample := range samples {
					// generate random UUID with ZBProxy signature
					fastrand.Read(u[:])
					u[0] = byte(i)
					u[1] = '$'
					u[2] = 'Z'
					u[3] = 'B'
					u[4] = '$'

					// marshal UUID string
					const hexTable = "0123456789abcdef"
					dst[8] = '-'
					dst[13] = '-'
					dst[18] = '-'
					dst[23] = '-'
					for i, x := range [16]byte{
						0, 2, 4, 6,
						9, 11,
						14, 16,
						19, 21,
						24, 26, 28, 30, 32, 34,
					} {
						c := u[i]
						dst[x] = hexTable[c>>4]
						dst[x+1] = hexTable[c&0x0F]
					}

					convertedSamples = append(convertedSamples, Sample{
						Name: sample.(string),
						ID:   string(dst[:]),
					})
				}

			default:
				log.Println(color.HiMagentaString(
					"Config Reload : Failed to reload samples: unknown samples input type: %T", samples))
				return false
			}
			s.Minecraft.OnlineCount.Sample = convertedSamples
		}
	}

	Config = config
	debug.FreeOSMemory()
	return true
}
