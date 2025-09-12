  import { Injectable } from '@angular/core';
  import { HttpClient, HttpHeaders } from '@angular/common/http';
  import { Observable } from 'rxjs';
  import { ConfigService } from '../../../services/config.service';
  
  @Injectable({
    providedIn: 'root',
  })
  export class AirLineService {//Carregar todos os diários sem documentos


    
  constructor(
    private http: HttpClient,
    private configService:ConfigService
  ) {}

    getAllAirLine(): Observable<any> {
      return this.http.get(this.configService.apiUrl + "/GetAllAirline" , {
        headers: new HttpHeaders({
          'Content-Type': 'application/json',
        }),
      });
    }
  

}