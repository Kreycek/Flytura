import { Injectable } from '@angular/core';

import { TokenPayload } from '../interfaces/tokenPayload';
import { jwtDecode } from 'jwt-decode';
import { TranslateService } from '@ngx-translate/core';

@Injectable({
  providedIn: 'root',
})
export class ConfigService {

    // public apiUrl = 'http://localhost:8080'; 
    public apiUrl = ''; 

    constructor() {
      

        const hostname = window.location.hostname;
        const protocol = window.location.protocol;
        console.log('hostname',hostname);
        console.log('protocol',protocol)

        if (hostname.includes('54.156.244.197') && (protocol=="https:" || protocol=="http:")) {
            //AMBIENTE DE PRODUÇÃO
            this.apiUrl = protocol + '//54.156.244.197/api';
        }      
        else if(hostname.includes('18.210.18.180') &&  protocol=="http:") {
            //AMBIENTE DE HOMOL
            this.apiUrl = protocol + '//18.210.18.180/api';
        }  
        else if((hostname.includes('app.flytura.com')) &&  protocol=="https:") {
             this.apiUrl = protocol + 'api';
        }   
        else if((hostname.includes('127.0.0.1:4200')) &&  protocol=="http:") {
             this.apiUrl = protocol + '//localhost:8080/api';
        }        
        else {
            this.apiUrl = protocol + '//localhost:8080/api';
        }
       

        // if (hostname.includes('54.156.244.197') && protocol=="https:") {
        //     this.apiUrl = 'http://54.156.244.197:8080';
        // }
        // else if(hostname.includes('18.210.18.180')) {
        //     this.apiUrl = 'http://18.210.18.180:8080';
        // } else {
        //     this.apiUrl = 'http://localhost:8080';
        // }
    }


    
    public limitPaginator=50;

    
   isTokenExpired(expSeconds: number, clockSkewSeconds = 30): boolean {
        const nowSeconds = Math.floor(Date.now() / 1000);
        return nowSeconds >= (expSeconds - clockSkewSeconds);
    }

    tokenExpired() :boolean  {
        const token = localStorage.getItem('token');
        if(token) {
            const payload =  jwtDecode<TokenPayload>(token)

            if (payload && this.isTokenExpired(payload.exp??0)) {
                    localStorage.removeItem('token'); 
                    return true;
            }
            else {
                return false;
            }
        } else {
            return true;
        }
    }

    returnTokenData() : TokenPayload | null {
        
        const token = localStorage.getItem('token'); // ou onde você armazenou o JWT

        if (token) {
        return  jwtDecode<TokenPayload>(token); 
        }
        else {
        return null
        }
    }
    


    public months= [
        { "value": 1, "name": "Janeiro", "shortName": "Jan" },
        { "value": 2, "name": "Fevereiro", "shortName": "Fev" },
        { "value": 3, "name": "Março", "shortName": "Mar" },
        { "value": 4, "name": "Abril", "shortName": "Abr" },
        { "value": 5, "name": "Maio", "shortName": "Mai" },
        { "value": 6, "name": "Junho", "shortName": "Jun" },
        { "value": 7, "name": "Julho", "shortName": "Jul" },
        { "value": 8, "name": "Agosto", "shortName": "Ago" },
        { "value": 9, "name": "Setembro", "shortName": "Set" },
        { "value": 10, "name": "Outubro", "shortName": "Out" },
        { "value": 11, "name": "Novembro", "shortName": "Nov" },
        { "value": 12, "name": "Dezembro", "shortName": "Dez" }
      ]

/*
Função criada por Ricardo Silva Ferreira
Início da criação: 24/03/2026 15:11
Data final da criação: 24/03/2026 15:11
*/
    sortByKey(array:any, key:any) {
        return array.slice().sort((a:any, b:any) => {
            const valA = String(a[key]).toLowerCase();
            const valB = String(b[key]).toLowerCase();

            if (valA < valB) return -1;
            if (valA > valB) return 1;
            return 0;
        });
    }


}