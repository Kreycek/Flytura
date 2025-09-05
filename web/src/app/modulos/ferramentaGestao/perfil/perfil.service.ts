import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { catchError } from 'rxjs/operators';
import { ConfigService } from '../../../services/config.service';



@Injectable({
  providedIn: 'root',
})
export class PerfilService {
    // URL do seu backend para login

  constructor(
    private http: HttpClient,
    private configService:ConfigService
  ) {}

  gePerfil(): Observable<any> {
    return this.http.get(this.configService.apiUrl + "/getPerfis", {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      }),
    });
  }
}
