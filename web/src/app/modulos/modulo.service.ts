
  import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { FormArray, FormGroup, Validators } from '@angular/forms';
import { Observable } from 'rxjs';
import { ConfigService } from '../services/config.service';
  
  
  
  
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


     


}