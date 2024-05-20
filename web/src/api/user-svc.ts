export interface User {
    firstName: string
    lastName: string
    name: string
    roles: string[]
}

export async function authenticate(username: string, password: string): Promise<User | null> {
    let user: User | null = null
    try {
        const resp = await fetch("/api/login", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ username, password })
        })

        if (!resp.ok) {
            console.warn(`There was a problem attempting to login. Server responded with code ${resp.status}: ${resp.statusText}`)
        } else {
            user = await resp.json() as User
        }
    } catch (error) {
        console.warn(`Login of user ${username} failed: ${error}`)
    }

    return user
}
