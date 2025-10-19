import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { ConfigService } from '../../../services/config.service';

@Injectable({
  providedIn: 'root',
})
export class InvoicesService {
    // URL do seu backend para login

  constructor(
    private http: HttpClient,
    private configService:ConfigService
  ) {}


  
  getAllS3ImagesDBlDataPagination(page:number, limit:number, companyCode?:string | null, startDate?:string | null, endDate?:string | null): Observable<any> {
    return this.http.get(
        this.configService.apiUrl + "/SearchS3Images?page="+page + 
        "&limit="+limit + 
        "&companyCode="+companyCode+ 
        "&startDate="+startDate+ 
        "&endDate="+endDate , {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      }),
    });
  }

}