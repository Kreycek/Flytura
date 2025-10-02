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

  getAllOnlyFlyExcelDataPagination(page:number, limit:number): Observable<any> {
    return this.http.get(this.configService.apiUrl + "/GetAllOnlyFlyExcelData?page="+page + "&limit="+limit , {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      }),
    });
  }

    //Carregar todos os diários sem documentos
    getAllOnlyFlyExcelData(): Observable<any> {
      return this.http.get(this.configService.apiUrl + "/GetOnlyFlyExcelData" , {
        headers: new HttpHeaders({
          'Content-Type': 'application/json',
        }),
      });
    }
  

  getOnlyFlyExcelDataById(id:string): Observable<any> {
    return this.http.get(this.configService.apiUrl + "/GetOnlyFlyExcelDataById?id="+id, {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      }),
    });
  }

  addOnlyFlyExcelData(data:any): Observable<any> {
    return this.http.post(this.configService.apiUrl + "/InsertOnlyFlyExcelData", data, {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      }),
    });
  }

  updateOnlyFlyExcelData(data:any): Observable<any> {
    return this.http.post(this.configService.apiUrl + "/UpdateOnlyFlyExcelData", data, {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      }),
    });
  }

  verifyExistOnlyFlyExcelData(data:any): Observable<any> {
    return this.http.post(this.configService.apiUrl + "/VerifyExistOnlyFlyExcelData", data, {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      }),
    });
  }

  searchOnlyFlyExcelData(data:any): Observable<any> {
    return this.http.post(this.configService.apiUrl + "/SearchOnlyFlyExcelData", data, {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      }),
    });
  }

  
  importarPlanilha(formData:FormData): Observable<any> {
    return this.http.post(this.configService.apiUrl + "/uploadOnlyFlyExcelData", formData);

  }


  //Carregar todos os status de importação
    getAllStatusImportData(): Observable<any> {
      return this.http.get(this.configService.apiUrl + "/GetAllImportStatus" , {
        headers: new HttpHeaders({
          'Content-Type': 'application/json',
        }),
      });
    }


    
  //Carregar todos os status de importação
   
GroupByCompanyName(
  status?: string | null,
  companyName?: string| null,
  startDate?: string| null,
  endDate?: string| null
): Observable<any> {
  let params = new HttpParams();

  if (status) {
    params = params.set('status', status);
  }
  if (companyName) {
    params = params.set('companyName', companyName);
  }
  if (startDate) {
    params = params.set('startDate', startDate);
  }
  if (endDate) {
    params = params.set('endDate', endDate);
  }

  return this.http.get(this.configService.apiUrl + '/GroupByCompanyName', {
    headers: new HttpHeaders({
      'Content-Type': 'application/json',
    }),
    params: params
  });
}

}
