import { HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { Router } from '@angular/router';
import { HttpRequest, HttpHandlerFn, HttpErrorResponse } from '@angular/common/http';
import { catchError, tap, throwError } from 'rxjs';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { ConfigService } from '../services/config.service';



export const AuthInterceptor: HttpInterceptorFn = (req: HttpRequest<any>, next: HttpHandlerFn) => {

 

   const translateService = inject(TranslateService); 



    // Definir URLs que devem ser ignoradas pelo interceptor
    const ignoredUrls = ['/login']; // Adicione outras se necessário

    // Se a URL fizer parte das ignoradas, não modificar a requisição
    if (ignoredUrls.some(url => req.url.includes(url))) {
        return next(req);
    }

    // if(config.tokenExpired()) {
    //     setTimeout(() => {
    //              const titulo = translateService.instant('expirou');
    //           if (confirm(titulo)) {     
    //     //   window.location.assign('https://flytura.com');    
    //           }    
    //     },5000);  
    // }


        // Obtém o token do localStorage
    const token = localStorage.getItem('token');
    
    // Se o token existir, adiciona no cabeçalho da requisição
    if (token) {
        // Clona a requisição e adiciona o token no cabeçalho Authorization
        req = req.clone({
        setHeaders: {
            Authorization: `Bearer ${token}`
        }
        });
    }

    
 // Se você guarda o token inteiro, pode checar aqui antes de enviar
  

  var qtdViewMsg=0

  return next(req).pipe(
   tap(),
    catchError((error: HttpErrorResponse) => {
      
        if (error.status === 401 || error.status === 403) {
         
            const titulo = translateService.instant('Ecra.lostSesssion');
            if (confirm(titulo)) {     
                if(qtdViewMsg==0) {           
                    window.location.assign('https://flytura.com'); 
                    qtdViewMsg=1  
                }        
            } else {
                 if(qtdViewMsg==0) {           
                    window.location.assign('https://flytura.com');   
                      qtdViewMsg=1  
                }  
            }

              


            // window.location.assign('https://flytura.com');
            
            // router.navigate(['/login']); // Redireciona para a página de login
          }
      return throwError(() => error);
    })
  );
};
