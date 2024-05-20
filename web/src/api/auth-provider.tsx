import {createContext, FC, ReactElement, useState} from "react"
import {authenticate, User} from './user-svc'

type authSuccessHandler = (user: User) => void

export interface AuthContextValue {
    user: User | null
    isAuthenticated: () => boolean
    isAuthenticating: boolean
    login: (username: string, password: string, onAuthSuccess: authSuccessHandler) => void
}

const AuthContext = createContext<AuthContextValue>({
    user: null,
    isAuthenticated: () => false,
    isAuthenticating: false,
    login: (username: string, password: string, onAuthSuccess: authSuccessHandler) => {
    }
})

const AuthProvider: FC<{ children: ReactElement }> = ({ children }) => {

    const [user, setUser] = useState<User | null>(null)
    const [isAuthenticating, setIsAuthenticating] = useState(false)
    const isAuthenticated = () => user == null

    const login = async (username: string, password: string, onAuthSuccess = (user: User) => {
    }) => {
        setIsAuthenticating(true)
        const maybeUser = await authenticate(username, password)
        setUser(maybeUser)
        if (maybeUser != null) onAuthSuccess(maybeUser)
        setIsAuthenticating(false)
    }

    return <AuthContext.Provider value={{ user, isAuthenticated, isAuthenticating, login }}>{children}</AuthContext.Provider>
}

export { AuthContext, AuthProvider }