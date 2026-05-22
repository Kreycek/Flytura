import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { ConfigService } from '../../../../services/config.service';
import JSZip from 'jszip';
import { saveAs } from 'file-saver'

@Injectable({
  providedIn: 'root',
})
export class InvoicesService {
    // URL do seu backend para login

  constructor(
    private http: HttpClient,
    private configService:ConfigService
  ) {}
  
  getAllS3ImagesDBDataPagination(page:number, limit:number,billedFlytura:boolean,doDownload:boolean, key?:string | null, companyCode?:string | null, startDate?:string | null, endDate?:string | null): Observable<any> {
     console.log('this.objPesquisar.endDate',endDate)

    return this.http.get(
        this.configService.apiUrl + "/SearchS3ImagesDBPagination?page="+page + 
        "&limit="+limit + 
         "&billedFlytura="+billedFlytura+ 
         "&doDownload="+doDownload+ 
         "&key="+key+ 
        "&companyCode="+companyCode+ 
        "&startDate="+startDate+ 
        "&endDate="+(endDate && endDate!=undefined ? endDate : '')  , {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      }),
    });
  }

 
  getAllS3ImagesDBFull(companyCode?:string | null, startDate?:string | null, endDate?:string | null): Observable<any> {
    return this.http.get(
        this.configService.apiUrl + "/SearchS3ImagesDBFull?companyCode="+companyCode+ 
        "&startDate="+startDate+ 
        "&endDate="+endDate , {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      }),
    });
  }


  async downloadZip(urls: string[]) {
    const zip = new JSZip();

    for (const url of urls) {
      try {
        const response = await fetch(url);
        const blob = await response.blob();
        const filename = url.split('/').pop() || 'file';
        zip.file(filename, blob);
      } catch (error) {
        console.error(`Erro ao baixar ${url}:`, error);
      }
    }

    zip.generateAsync({ type: 'blob' }).then((content) => {
      saveAs(content, 'Facturas.zip');
    });
  }


    /*
      Função criada por Ricardo Silva Ferreira
      Inicio da criação 02/12/2025 14:54
      Data Final da criação : 02/12/2025 14:54
      Obs: Agrupa varios arquivos zipados por compania aerea, ou seja cria um arquivo fatcuras.zip e dentro outros zips
          separados por companhia aérea
    */
  
async downloadGroupedZip(items: any[]) {
  const mainZip = new JSZip();

  // Agrupar por CompanyName
  const grouped: Record<string, string[]> = items.reduce((acc, item) => {
    if (!acc[item.CompanyName]) acc[item.CompanyName] = [];
    acc[item.CompanyName].push(item.ZipFileName);
    return acc;
  }, {} as Record<string, string[]>);

  // Para cada empresa, criar uma pasta dentro do ZIP principal
  for (const [companyName, urls] of Object.entries(grouped) as [string, string[]][]) {
    const folder = mainZip.folder(companyName);
    if (!folder) continue;

    for (const url of urls) {
      try {
        const response = await fetch(url);
        const blob = await response.blob();
        const filename = url.split('/').pop() || 'file.zip';
        folder.file(filename, blob);
      } catch (error) {
        console.error(`Erro ao baixar ${url}:`, error);
      }
    }
  }

  // Gerar o ZIP final
  const content = await mainZip.generateAsync({ type: 'blob' });
    saveAs(content, 'Facturas.zip');
  }

  async downloadGroupedZipAndReturnTrue(items: any): Promise<boolean> {
    await this.downloadGroupedZip(items);
    return true;
  }

  updateStatusS3Image(formData:any): Observable<any> {
    return this.http.post(this.configService.apiUrl + "/UpdateStatusS3Image", formData);
  }

  updateMultipleStatusS3Images(formData:any): Observable<any> {
    return this.http.post(this.configService.apiUrl + "/UpdateMultipleStatusS3Images", formData);
  }

  UpdateStatusPdforXml(formData:any): Observable<any> {
    return this.http.post(this.configService.apiUrl + "/UpdateStatusPdforXml", formData);
  }
  
  deleteS3Images(id:string,zipFile:string,pdfFile:string,xmlFile:string): Observable<any> {
    return this.http.delete(this.configService.apiUrl + '/DeleteImagesDBByIDHandler?id='+ id + 
      "&zipFile="+ zipFile+ 
      "&pdfFile="+pdfFile+ 
      "&xmlFile="+xmlFile, {
     
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      }),
    });
  }
  
  
}
