import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { ConfigService } from '../services/config.service';

@Injectable({
  providedIn: 'root',
})
export class AuthService {

 

  constructor(private http: HttpClient, private configService: ConfigService) {}

  validateToken(): Observable<any> {
    // Recuperar o token do localStorage
    const token = localStorage.getItem('token');

    if (!token || token==undefined) {

      return new Observable((observer) => {
        observer.error('Token não encontrado');
      });
    }

    // Enviar o token no cabeçalho da requisição
    const headers = new HttpHeaders().set('Authorization', token);
    return this.http.get(this.configService.apiUrl + "/validate", { headers });
  }
}
