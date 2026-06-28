
  import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { FormArray, FormGroup, Validators } from '@angular/forms';
import { Observable } from 'rxjs';
import { ConfigService } from '../services/config.service';
import saveAs from 'file-saver';
import * as XLSX from 'xlsx';
  
  
  
  @Injectable({
    providedIn: 'root',
  })
  export class ModuloService {
  
      constructor(
        private http: HttpClient,
        private configService:ConfigService
      ) {}
  
    desabilitaCamposFormGroup(form:FormGroup) {
        Object.keys(form.controls).forEach(controlName => {

        const control = form.get(controlName);
        control?.clearValidators();
        control?.updateValueAndValidity();
        });
    }

    habilitaCamposFormGroup(form:FormGroup,camposAtivar:any[]) {
        Object.keys(form.controls).forEach(controlName => {
        const control = form.get(controlName);
        // Defina as validações conforme necessário

        const existe=camposAtivar.some((campo)=>campo===controlName)
        // if (controlName === 'codDocument' || controlName === 'description' || controlName === 'country') {
            if(existe)
            control?.setValidators([Validators.required]);
        // }
        // Adicione as validações conforme o caso
        control?.updateValueAndValidity();
        });
    }

    isFieldInvalid(fieldName: string,  form:any): boolean {


        const field = (form as FormGroup).get(fieldName);
      
        return !!(field && field.invalid && (field.dirty || field.touched));
      }

    forcarAtivarValidarores(fieldName: string, form:FormGroup): void {
        const field = form.get(fieldName);
        if (field) {
            field.markAsTouched(); // Marca como "tocado"
            field.markAsDirty();   // Marca como "modificado"
            field.updateValueAndValidity(); // Recalcula a validação
        }
    }

    ativarvalidadores(form:FormGroup) {
        Object.keys(form.controls).forEach(fieldName => {            
            this.forcarAtivarValidarores(fieldName,form);
        });
    }


    forcarDesativarValidarores(fieldName: string, form:FormGroup): void {
        const field = form.get(fieldName);
        if (field) {
            field.markAsPristine(); // Marca como "não modificado"
            field.markAsUntouched(); // Marca como "não tocado"
            field.updateValueAndValidity(); // Atualiza o estado da validação
        }
    }

    desativarValidadores(form:FormGroup) {
        Object.keys(form.controls).forEach(fieldName => {            
            this.forcarDesativarValidarores(fieldName,form);
        });
    }

    filterDocuments(codDaily:any, dailys:any[], addOptionAll:boolean) {
        
        let retorno:any[]=[]
          let _documents=[]
          _documents=dailys.filter((response:any)=>{
            return response.codDaily===codDaily
    
          } )[0]
    
          if(_documents) {
            retorno=[];
            
            if(_documents.documents && _documents.documents.length>0) {        
              retorno= [..._documents.documents];
              if(addOptionAll) {
              retorno.unshift({
                "codDocument": "",
                "description": "Todos",
                "dtAdd": ""
              } )
            }
            }
            else 
            retorno=[]    
          }
          else {
            retorno=[]
          }
    
          return retorno;
      }
      
      getLastDataCoin(daily:any): Observable<any> {
        return this.http.get("https://economia.awesomeapi.com.br/json/last/"+ daily, {
          headers: new HttpHeaders({
            'Content-Type': 'application/json',
          }),
        });
      }


      deleteFormArrayData(fa:FormArray) {
        while (fa.length !== 0) {
          fa.removeAt(0);
        }
        fa.reset();
      }    


  //Carregar todos os status de importação
    getAllStatusImportData(): Observable<any> {
      return this.http.get(this.configService.apiUrl + "/GetAllImportStatus" , {
        headers: new HttpHeaders({
          'Content-Type': 'application/json',
        }),
      });
    }


    /*
    Função criada por Ricardo Silva Ferreira
    Inicio da criação 25/11/2025 17:29
    Data Final da criação :  25/11/2025 17:29
    */
    getFileNameS3(filePath:string) {

      const parts = filePath.split('/');

      return parts[parts.length-1];
    }    

  /*
    Função criada por Ricardo Silva Ferreira
    Inicio da criação 09/04/2026 09:25
    Data Final da criação :  09/04/2026 09:30
    OBS converte a data da api para uma data de DD/MM/YYYY hh:mm
  */
  handleExcelData(letterColumns:string,ws:XLSX.WorkSheet, data:any[]) {
      for (let i = 2; i <= data.length + 1; i++) {
            const cell = ws[`${letterColumns}${i}`]; // coluna E = Import Date

            if (cell && cell.v) {
              // console.log('columns data',cell.v);
              const jsDate: Date = new Date(cell.v);

              // Converte JS Date -> Excel number date
              const excelNumber =
                (jsDate.getTime() - new Date(Date.UTC(1899, 11, 30)).getTime()) /
                86400000;

              cell.t = "n";                    // tipo number (Excel date)
              // cell.z = "dd/mm/yyyy hh:mm";     // formato desejado
              cell.z = "dd/mm/yyyy";     // formato desejado
              cell.v = excelNumber;            // número convertido
            }
          }
  }

  /*
  Função criada por Ricardo Silva Ferreira
  Inicio da criação 09/04/2026 09:25
  Data Final da criação :  09/04/2026 09:30
  */
  exportToExcelPurcharseRecord(apiData: any[], fileName: string): void {
        
      const dados = apiData.map(item => ({
            "Key": item.Key,
            "Name": item.Name,
            "EmissionDate": item.EmissionDate,
            "Last Name": item.LastName,
            "Airline Company": item.CompanyName,
            "Import Date": item.DtImportacao,
            "Status":item.Status,
            "Ida-Volta":item.DirectionOfDestination,
            "Msg Return Error": item.MessageReturn
          }));

          const colunas = ["Key", "Name","EmissionDate","Last Name", "Airline Company", "Import Date","Status","Ida-Volta",'Msg Return Error'];
          // Converte JSON para uma worksheet Excel
          const worksheet: XLSX.WorkSheet = XLSX.utils.json_to_sheet(dados, { header: colunas });

          // Transformamos a data em número Excel de forma segura

          this.handleExcelData('F',worksheet,dados)             

        // Cria a workbook
        const workbook: XLSX.WorkBook = {
          Sheets: { 'Dados': worksheet },
          SheetNames: ['Dados']
        };

        // Converte workbook para buffer
        const excelBuffer: any = XLSX.write(workbook, {
          bookType: 'xlsx',
          type: 'array'
        });

        // Salva arquivo
        const blob: Blob = new Blob([excelBuffer], {
          type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet;charset=UTF-8'
        });

    saveAs(blob, `${fileName}.xlsx`);
  }

  /*
  Função criada por Ricardo Silva Ferreira
  Inicio da criação 09/04/2026 10:15
  Data Final da criação :  09/04/2026 10:15
  */
  exportToExcelInformes(apiData: any[], fileName: string): void {
  
    const dados = apiData.map(item => ({
            "Key": item.Key,
            "Status": item.Status,
            "Processing Date": item.DtProcess,
            "Processing Month": item.MonthProcess,            
            "Flight date":item.DtFly,
            "RFC": item.RFC,
            "T. Base": item.TransferredBaseValue,
            "IVA": item.IVAValue,
            "TAX": item.Tax,
            'Rate':item.Rate,
            'Factor Type':item.FactorType,           
            'Sub Total':item.SubTotalValue,
            'Total':item.TotalValue,
            'Ruta':item.Ruta,
            'Airline Company':item.CompanyName,
            'Tua':item.TUA,
            'Other values':item.OtherValues,
            'Segment':item.Segment,
            'Ticket':item.Ticket
          }));

          const colunas = [
            "Key", 
            "Status", 
            "Processing Date", 
            "Processing Month", 
            "Flight date",
            'RFC',
            'T. Base',
            'IVA',
            'TAX',
            'Rate',
            'Factor Type',
            'Sub Total',
            'Total',
            'Ruta',
            'Airline Company',
            'Tua',
            'Other values',
            'Segment'];
          // Converte JSON para uma worksheet Excel
          const worksheet: XLSX.WorkSheet = XLSX.utils.json_to_sheet(dados, { header: colunas });

          // Transformamos a data em número Excel de forma segura
          this.handleExcelData('C',worksheet,dados)             
          

        // Cria a workbook
        const workbook: XLSX.WorkBook = {
          Sheets: { 'Dados': worksheet },
          SheetNames: ['Dados']
        };

        // Converte workbook para buffer
        const excelBuffer: any = XLSX.write(workbook, {
          bookType: 'xlsx',
          type: 'array'
        });

        // Salva arquivo
        const blob: Blob = new Blob([excelBuffer], {
          type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet;charset=UTF-8'
        });

    saveAs(blob, `${fileName}.xlsx`);
  }

  
   /*
  Função criada por Ricardo Silva Ferreira
  Inicio da criação 22/05/2026 14:46
  Data Final da criação :  22/05/2026 14:46
  */
  public isValidDate(date: any): boolean {
      if (!date) return false;

      const d = new Date(date);

      return d.getFullYear() > 1; // ignora 0001
  }



  
  /*
  Função criada por Ricardo Silva Ferreira
  Inicio da criação 22/05/2026 16:57
  Data Final da criação :  22/05/2026 10:15
  */
  exportToExcelConciliation(apiData: any[], fileName: string): void {
  
    const dados = apiData.map(item => ({
            "Protocol": item.Protocol,
            "EmissionDate":  (item.EmissionDate && item.EmissionDate!='0001-01-01T00:00:00Z') ? item.EmissionDate : '',
            "OriginDate": (item.OriginDate && item.ReturnDate!='0001-01-01T00:00:00Z') ? item.OriginDate : '',
            "ReturnDate": (item.ReturnDate && item.ReturnDate!='0001-01-01T00:00:00Z') ? item.ReturnDate : '',
            "OriginLocator": item.OriginLocator,            
            "ReturnLocator":item.ReturnLocator,
            "OriginETicket": item.OriginETicket,
            "ReturnETicket": item.ReturnETicket,
            "OriginAirline": item.OriginAirline,
            "ReturnAirline": item.ReturnAirline,
            'TravelerName':item.TravelerName,
            'TravelerFirstName':item.TravelerFirstName,           
            'TravelerLastName':item.TravelerLastName,
            'BookingStatus':item.BookingStatus,
            'BookingStatusOld':item.BookingStatusOld,
            'OriginCancelledReason':item.OriginCancelledReason,
            'ReturnCancelledReason':item.ReturnCancelledReason,
            'CurrencyCode':item.CurrencyCode,
            'AmountOrigin':item.AmountOrigin,
            'AmountReturn':item.AmountReturn,
            'CreatedAtLocalCountry':item.CreatedAtLocalCountry,
            'CreatedAtContractedCountry':item.CreatedAtContractedCountry,
            "FlightOrigin":item.FlightOrigin,
            "FlightDestination":item.FlightDestination,
            "OriginCountryCode":item.OriginCountryCode,
            "OriginCity":item.OriginCity,
            "DestinationCountryCode":item.DestinationCountryCode,
            "DestinationCity":item.DestinationCity,
            "OriginAirlineCommercial":item.OriginAirlineCommercial,
            "ReturnAirlineCommercial":item.ReturnAirlineCommercial,
            "MxnOnflyAmountOrigin":item.MxnOnflyAmountOrigin,
            "MxnOnflyAmountReturn":item.MxnOnflyAmountReturn,

          }));

          const colunas = [
            "Protocol", 
            "EmissionDate", 
            "OriginDate",
            "ReturnDate", 
            "OriginLocator", 
            "ReturnLocator",
            'OriginETicket',
            'ReturnETicket',
            'OriginAirline',
            'ReturnAirline',
            'TravelerName',
            'TravelerFirstName',
            'TravelerLastName',
            'BookingStatus',
            'BookingStatusOld',
            'OriginCancelledReason',
            'ReturnCancelledReason',
            'CurrencyCode',
            'AmountOrigin',
            'AmountReturn',
            'CreatedAtLocalCountry',
            'CreatedAtContractedCountry',
            "FlightOrigin",
            "FlightDestination",
            "OriginCountryCode",
            "OriginCity",
            "DestinationCountryCode",
            "DestinationCity",
            "OriginAirlineCommercial",
            "ReturnAirlineCommercial",
            "MxnOnflyAmountOrigin",
            "MxnOnflyAmountReturn"
          ];
          // Converte JSON para uma worksheet Excel
          const worksheet: XLSX.WorkSheet = XLSX.utils.json_to_sheet(dados, { header: colunas });

          // Transformamos a data em número Excel de forma segura
          this.handleExcelData('C',worksheet,dados)             
          

        // Cria a workbook
        const workbook: XLSX.WorkBook = {
          Sheets: { 'Dados': worksheet },
          SheetNames: ['Dados']
        };

        // Converte workbook para buffer
        const excelBuffer: any = XLSX.write(workbook, {
          bookType: 'xlsx',
          type: 'array'
        });

        // Salva arquivo
        const blob: Blob = new Blob([excelBuffer], {
          type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet;charset=UTF-8'
        });

    saveAs(blob, `${fileName}.xlsx`);
  }

  
  /*
  Função criada por Ricardo Silva Ferreira
  Inicio da criação 01/01/2026 09:09
  Data Final da criação :  01/01/2026 09:09
  */
   importFile(formData:FormData, apiName:string): Observable<any> {
    return this.http.post(this.configService.apiUrl + apiName, formData);

  }
}