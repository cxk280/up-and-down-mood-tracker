package models

import (
  "errors"
  "time"

  "golang-mood-tracker/db"
  "golang-mood-tracker/forms"
)

//dashboard ...
type Dashboard struct {
  ID        int64    `db:"id, primarykey, autoincrement" json:"id"`
  UserID    int64    `db:"user_id" json:"-"`
  Title     string   `db:"title" json:"title"`
  Content   string   `db:"content" json:"content"`
  UpdatedAt int64    `db:"updated_at" json:"updated_at"`
  CreatedAt int64    `db:"created_at" json:"created_at"`
  User      *JSONRaw `db:"user" json:"user"`
}

//DashboardModel ...
type DashboardModel struct{}

//Create ...
func (m DashboardModel) Create(userID int64, form forms.DashboardForm) (dashboardID int64, err error) {
  getDb := db.GetDB()

  userModel := new(UserModel)

  checkUser, err := userModel.One(userID)

  if err != nil && checkUser.ID > 0 {
    return 0, errors.New("User doesn't exist")
  }

  _, err = getDb.Exec("INSERT INTO public.dashboard(user_id, title, content, updated_at, created_at) VALUES($1, $2, $3, $4, $5) RETURNING id", userID, form.Title, form.Content, time.Now().Unix(), time.Now().Unix())

  if err != nil {
    return 0, err
  }

  dashboardID, err = getDb.SelectInt("SELECT id FROM public.dashboard WHERE user_id=$1 ORDER BY id DESC LIMIT 1", userID)

  return dashboardID, err
}

//One ...
func (m DashboardModel) One(userID, id int64) (dashboard Dashboard, err error) {
  err = db.GetDB().SelectOne(&dashboard, "SELECT a.id, a.title, a.content, a.updated_at, a.created_at, json_build_object('id', u.id, 'name', u.name, 'email', u.email) AS user FROM public.dashboard a LEFT JOIN public.user u ON a.user_id = u.id WHERE a.user_id=$1 AND a.id=$2 LIMIT 1", userID, id)
  return dashboard, err
}

//All ...
func (m DashboardModel) All(userID int64) (dashboards []Dashboard, err error) {
  _, err = db.GetDB().Select(&dashboards, "SELECT a.id, a.title, a.content, a.updated_at, a.created_at, json_build_object('id', u.id, 'name', u.name, 'email', u.email) AS user FROM public.dashboard a LEFT JOIN public.user u ON a.user_id = u.id WHERE a.user_id=$1 ORDER BY a.id DESC", userID)
  return dashboards, err
}

//Update ...
func (m DashboardModel) Update(userID int64, id int64, form forms.DashboardForm) (err error) {
  _, err = m.One(userID, id)

  if err != nil {
    return errors.New("dashboard not found")
  }

  _, err = db.GetDB().Exec("UPDATE public.dashboard SET title=$1, content=$2, updated_at=$3 WHERE id=$4", form.Title, form.Content, time.Now().Unix(), id)

  return err
}

//Delete ...
func (m DashboardModel) Delete(userID, id int64) (err error) {
  _, err = m.One(userID, id)

  if err != nil {
    return errors.New("dashboard not found")
  }

  _, err = db.GetDB().Exec("DELETE FROM public.dashboard WHERE id=$1", id)

  return err
}