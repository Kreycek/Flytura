import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class AuthService {

  private apiUrl = 'http://localhost:8080/validate';  // URL do seu backend para validar token

  constructor(private http: HttpClient) {}

  validateToken(): Observable<any> {
    // Recuperar o token do localStorage
    const token = localStorage.getItem('token');
    console.log('token ',token);
    if (!token || token==undefined) {
      console.log('token encontrado',token);
      return new Observable((observer) => {
        observer.error('Token não encontrado');
      });
    }

    // Enviar o token no cabeçalho da requisição
    const headers = new HttpHeaders().set('Authorization', token);
    return this.http.get(this.apiUrl, { headers });
  }
}
