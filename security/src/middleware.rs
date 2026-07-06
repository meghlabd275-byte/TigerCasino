use actix_web::{
    dev::{ServiceRequest, ServiceResponse, Transform, Service},
    Error, HttpMessage, HttpRequest,
    body::BoxBody,
};
use std::future::{Ready, Future};
use std::pin::Pin;
use std::rc::Rc;
use std::cell::RefCell;
use uuid::Uuid;

// ============== Auth Middleware ==============

pub struct AuthMiddleware;

impl<S, B> Transform<S, ServiceRequest> for AuthMiddleware
where
    S: Service<ServiceRequest, Response = ServiceResponse<B>, Error = Error> + 'static,
    S::Future: 'static,
    B: 'static,
{
    type Response = ServiceResponse<BoxBody>;
    type Error = Error;
    type InitError = ();
    type Transform = AuthMiddlewareService<S>;
    type Future = Ready<Result<Self::Transform, Self::InitError>>;

    fn new_transform(&self, service: S) -> Self::Future {
        ready(Ok(AuthMiddlewareService {
            service: Rc::new(RefCell::new(service)),
        }))
    }
}

pub struct AuthMiddlewareService<S> {
    service: Rc<RefCell<S>>,
}

impl<S, B> Service<ServiceRequest> for AuthMiddlewareService<S>
where
    S: Service<ServiceRequest, Response = ServiceResponse<B>, Error = Error> + 'static,
    S::Future: 'static,
    B: 'static,
{
    type Response = ServiceResponse<BoxBody>;
    type Error = Error;
    type Future = Pin<Box<dyn Future<Output = Result<Self::Response, Self::Error>>>>;

    fn poll_ready(&self, cx: &mut std::task::Context<'_>) -> std::task::Poll<Result<(), Self::Error>> {
        self.service.borrow_mut().poll_ready(cx)
    }

    fn call(&self, req: ServiceRequest) -> Self::Future {
        let service = Rc::clone(&self.service);

        Box::pin(async move {
            // Check for authorization header
            let auth_header = req.headers()
                .get("authorization")
                .and_then(|v| v.to_str().ok());

            if let Some(auth) = auth_header {
                if auth.starts_with("Bearer ") {
                    let token = &auth[7..];
                    
                    // Simple token parsing (in production, validate JWT properly)
                    if let Some(user_id_str) = token.split('.').next() {
                        if let Ok(user_id) = Uuid::parse_str(user_id_str) {
                            // Store user_id in request extensions
                            req.extensions_mut().insert(user_id);
                        }
                    }
                }
            }

            let res = service.borrow_mut().call(req).await?;
            Ok(res)
        })
    }
}

// ============== Admin Middleware ==============

pub struct AdminMiddleware;

impl<S, B> Transform<S, ServiceRequest> for AdminMiddleware
where
    S: Service<ServiceRequest, Response = ServiceResponse<B>, Error = Error> + 'static,
    S::Future: 'static,
    B: 'static,
{
    type Response = ServiceResponse<BoxBody>;
    type Error = Error;
    type InitError = ();
    type Transform = AdminMiddlewareService<S>;
    type Future = Ready<Result<Self::Transform, Self::InitError>>;

    fn new_transform(&self, service: S) -> Self::Future {
        ready(Ok(AdminMiddlewareService {
            service: Rc::new(RefCell::new(service)),
        }))
    }
}

pub struct AdminMiddlewareService<S> {
    service: Rc<RefCell<S>>,
}

impl<S, B> Service<ServiceRequest> for AdminMiddlewareService<S>
where
    S: Service<ServiceRequest, Response = ServiceResponse<B>, Error = Error> + 'static,
    S::Future: 'static,
    B: 'static,
{
    type Response = ServiceResponse<BoxBody>;
    type Error = Error;
    type Future = Pin<Box<dyn Future<Output = Result<Self::Response, Self::Error>>>>;

    fn poll_ready(&self, cx: &mut std::task::Context<'_>) -> std::task::Poll<Result<(), Self::Error>> {
        self.service.borrow_mut().poll_ready(cx)
    }

    fn call(&self, req: ServiceRequest) -> Self::Future {
        let service = Rc::clone(&self.service);

        Box::pin(async move {
            // Check for admin flag in extensions (set by auth middleware)
            // In production, validate admin status from database
            let is_admin = req.extensions()
                .get::<bool>()
                .copied()
                .unwrap_or(false);

            // For now, allow all authenticated requests
            // In production, check admin role from database

            let res = service.borrow_mut().call(req).await?;
            Ok(res)
        })
    }
}

// ============== Helper Functions ==============

/// Extract user ID from request (must be used after auth middleware)
pub fn get_user_id(req: &HttpRequest) -> Option<Uuid> {
    req.extensions().get::<Uuid>().copied()
}

/// Check if request is from admin (must be used after admin middleware)
pub fn is_admin(req: &HttpRequest) -> bool {
    req.extensions().get::<bool>().copied().unwrap_or(false)
}
