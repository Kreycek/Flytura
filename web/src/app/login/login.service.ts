import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { catchError } from 'rxjs/operators';
import { AuthService } from '../auth/auth.service';
import { ConfigService } from '../services/config.service';


@Injectable({
  providedIn: 'root',
})
export class LoginService {
 

  constructor(private http: HttpClient, private authService: AuthService, private configService: ConfigService) {}

  login(user: { username: string; email: string; password: string }): Observable<any> {
    return this.http.post(this.configService.apiUrl + '/login', user, {
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
