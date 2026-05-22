import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { catchError } from 'rxjs/operators';
import { ConfigService } from '../../../services/config.service';



@Injectable({
  providedIn: 'root',
})
export class UsuarioService {
    // URL do seu backend para login

  constructor(
    private http: HttpClient,
    private configService:ConfigService
  ) {}

  getUsers(page:number, limit:number): Observable<any> {
    return this.http.get(this.configService.apiUrl + "/getAllUsers?page="+page + "&limit="+limit , {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      }),
    });
  }

  
  addUsers(user:any): Observable<any> {
    return this.http.post(this.configService.apiUrl + "/addUser", user, {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      }),
    });
  }

  updateUser(user:any): Observable<any> {
    return this.http.post(this.configService.apiUrl + "/updateUser", user, {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      }),
    });
  }


  verifyExistsUsers(user:any): Observable<any> {
    return this.http.post(this.configService.apiUrl + "/verifyExistUser", user, {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      }),
    });
  }

  searchUsers(user:any): Observable<any> {
    return this.http.post(this.configService.apiUrl + "/searchUsers", user, {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      }),
    });
  }

  getUserById(id:string): Observable<any> {
    return this.http.get(this.configService.apiUrl + "/getUserById?id="+id, {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      }),
    });
  }
}
