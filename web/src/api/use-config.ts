import {useEffect, useState} from "react";


export interface Config {
    title: string
}

const defaultConfig: Config = {
    title: ""
}

const useConfig = (): Config => {
    const [config, setConfig] = useState<Config>(defaultConfig)

    useEffect(  () => {

        async function doIt() {
            try {
                const resp = await fetch("/config.json")
                const json = await resp.json() as Config
                setConfig(json)
            } catch (err) {
                console.warn(`failed to retrieve config from backend: ${err}. Proceeding with default values:`, defaultConfig)
            }
        }

        doIt()

    }, [])


    return config

}

export default useConfig
