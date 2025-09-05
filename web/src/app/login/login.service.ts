import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { catchError } from 'rxjs/operators';
import { AuthService } from '../auth/auth.service';


@Injectable({
  providedIn: 'root',
})
export class LoginService {
  private apiUrl = 'http://localhost:8080/login';  // URL do seu backend para login

  constructor(private http: HttpClient, private authService: AuthService) {}

  login(user: { username: string; email: string; password: string }): Observable<any> {
    return this.http.post(this.apiUrl, user, {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      }),
    });
  }

  // Método para validar o token logo após o login
  validateTokenAfterLogin(token: string): Observable<any> {
    return this.authService.validateToken();
  }
}
