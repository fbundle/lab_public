from __future__ import annotations


class Function:
    def differentiate(self, v: Variable) -> Function:
        raise NotImplementedError

    def simplify(self) -> Function:
        return self


class Variable(Function):
    name: str

    def __init__(self, name: str):
        self.name = name

    def differentiate(self, v: Variable) -> Function:
        if v.name == self.name:
            return ConstantFunction(const=1)
        return ConstantFunction(const=0)

    def __str__(self):
        return self.name


class ConstantFunction(Function):
    const: float

    def __init__(self, const: float):
        self.const = const

    def differentiate(self, v: Variable) -> Function:
        return ConstantFunction(const=0)

    def __str__(self):
        return str(self.const)


class LinearFunction(Function):
    def __init__(self, coef: float, param: Function):
        self.coef = coef
        self.param = param

    def differentiate(self, v: Variable) -> Function:
        if isinstance(self.param, Variable) and self.param.name == v.name:
            return ConstantFunction(self.coef)
        return LinearFunction(coef=self.coef, param=self.param.differentiate(v=v))

    def __str__(self):
        return f"{self.coef}×{self.param}"

    def simplify(self):
        self.param = self.param.simplify()
        if self.coef == 1:
            return self.param


class SumFunction(Function):
    param_list: list[Function]

    def __init__(self, param_list: list[Function]):
        self.param_list = param_list

    def differentiate(self, v: Variable) -> Function:
        return SumFunction(param_list=[param.differentiate(v) for param in self.param_list])

    def __str__(self):
        return "(" + "+".join(map(str, self.param_list)) + ")"


class Prod2Function(Function):
    x: Function
    y: Function

    def __init__(self, x: Function, y: Function):
        self.x, self.y = x, y

    def differentiate(self, v: Variable) -> Function:
        return SumFunction(param_list=[
            Prod2Function(x=self.x.differentiate(v), y=self.y),
            Prod2Function(x=self.x, y=self.y.differentiate(v)),
        ])

    def __str__(self):
        return f"{self.x}×{self.y}"

    def simplify(self) -> Function:
        self.x = self.x.simplify()
        self.y = self.y.simplify()
        if isinstance(x, Variable) and isinstance(y, Variable) and x.name == y.name:
            return PowFunction(n=2, x=self.x)
        return self


class PowFunction(Function):
    n: int
    x: Function

    def __init__(self, n: int, x: Function):
        self.n = n
        self.x = x

    def differentiate(self, v: Variable) -> Function:
        return LinearFunction(coef=self.n, param=PowFunction(self.n - 1, x))

    def __str__(self):
        return f"({self.x})^({self.n})"


class ExpFunction(Function):
    x: Function

    def __init__(self, x: Function):
        self.x = x

    def differentiate(self, v: Variable) -> Function:
        return Prod2Function(x=self, y=self.x.differentiate(v))

    def __str__(self):
        return f"e^{self.x}"

    def simplify(self) -> Function:
        self.x = self.x.simplify()
        return self


# TODO : differentiate here already have chain rule and rule for diff multiple variables - we want to bake chain rule into our code
# TODO : chain rule: D(f g)_x = D(f)_{g(x)} D(g)_x
# TODO : multivariate: f(t) = f(x(t), y(t)) then df/dt = (∂f/∂x) (dx/dt) + (∂f/∂y) (dy/dt)
# TODO : that is, user just need to specify partial derivative for each param


if __name__ == "__main__":
    x = Variable("x")
    y = Variable("y")
    a1 = LinearFunction(coef=2, param=x)
    a2 = SumFunction(param_list=[x, y, a1])
    a3 = Prod2Function(x=x, y=SumFunction([x, y]))
    a4 = LinearFunction(coef=3, param=a1)
    a5 = ExpFunction(x=Prod2Function(x=x, y=x))
    print(a5.simplify())
    print(a5.differentiate(x).simplify())
