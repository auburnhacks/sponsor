export interface User {
    id: string;
    name: string;
    email: string;
    ACL: string;
    token: string;
}

export interface Admin extends User {
}

export interface Sponsor extends User{
    company: Company;
}

export interface Company {
    id:   string;
    name: string;
    logo: string;
}