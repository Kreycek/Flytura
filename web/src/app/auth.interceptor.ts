import { HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { Router } from '@angular/router';
import { HttpRequest, HttpHandlerFn, HttpErrorResponse } from '@angular/common/http';
import { catchError, tap, throwError } from 'rxjs';

export const AuthInterceptor: HttpInterceptorFn = (req: HttpRequest<any>, next: HttpHandlerFn) => {
  const router = inject(Router); // Injeção manual do Router


    // Definir URLs que devem ser ignoradas pelo interceptor
    const ignoredUrls = ['/login']; // Adicione outras se necessário

    // Se a URL fizer parte das ignoradas, não modificar a requisição
    if (ignoredUrls.some(url => req.url.includes(url))) {
        return next(req);
    }


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


  return next(req).pipe(
   tap(),
    catchError((error: HttpErrorResponse) => {
      
        if (error.status === 401) {
            
            // router.navigate(['/login']); // Redireciona para a página de login
          }
      return throwError(() => error);
    })
  );
};
